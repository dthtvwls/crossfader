package crossfader

import (
	"bufio"
	"encoding/csv"
	"net"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

const (
	haproxySocket    = "/run/haproxy/admin.sock"
	haproxyConfFile  = "/etc/haproxy/haproxy.cfg"
	haproxyReloadCmd = "systemctl reload haproxy"
	haproxyConfTmpl  = `
    global
    	log /dev/log	local0
    	log /dev/log	local1 notice
    	chroot /var/lib/haproxy
    	stats socket /run/haproxy/admin.sock mode 660 level admin
    	stats timeout 30s
    	user haproxy
    	group haproxy
    	daemon

    defaults
    	log	 global
    	mode tcp
    	option	tcplog
    	timeout connect 5000
    	timeout client  50000
    	timeout server  50000

    resolvers mydns
    	nameserver dns1 8.8.8.8:53
    	nameserver dns2 8.8.4.4:53

    listen http-in
    	bind *:443{{range .}}
    	server {{.Server}} {{.Server}}:443 weight {{.Weight}} check resolvers mydns{{end}}
  `
)

func writeSocket(cmd string) (*net.UnixConn, error) {
	if socket, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: haproxySocket, Net: "unix"}); err != nil {
		return nil, err
	} else {
		writer := bufio.NewWriter(socket)
		if _, err := writer.WriteString(cmd + "\n"); err != nil {
			return nil, err
		} else if err := writer.Flush(); err != nil {
			return nil, err
		}
		return socket, nil
	}
}

func getHaproxyConf() (*Conf, error) {
	if socket, err := writeSocket("show stat -1 4 -1"); err != nil {
		return nil, err
	} else if records, err := csv.NewReader(bufio.NewReader(socket)).ReadAll(); err != nil {
		return nil, err
	} else {
		var conf Conf
		if len(records) > 2 {
			conf.Servers = [2]string{records[1][1], records[2][1]}
			if subtrahend, err := strconv.Atoi(records[2][18]); err != nil {
				return nil, err
			} else {
				conf.Subtrahend = subtrahend
			}
		}
		return &conf, nil
	}
}

func putHaproxyConf(value []byte) error {
	if conf, err := validate(value); err != nil {
		return err
	} else if haproxyConf, err := getHaproxyConf(); err != nil {
		return err
	} else if !reflect.DeepEqual(conf, haproxyConf) {
		// format data for template
		data := []struct {
			Server string
			Weight int
		}{
			{conf.Servers[0], 256 - conf.Subtrahend},
			{conf.Servers[1], conf.Subtrahend},
		}

		// rewrite config file
		if f, err := os.Create(haproxyConfFile); err != nil {
			return err
		} else if err := template.Must(template.New("").Parse(haproxyConfTmpl)).Execute(f, data); err != nil {
			return err
		}

		if conf.Servers[0] != haproxyConf.Servers[0] || conf.Servers[1] != haproxyConf.Servers[1] {
			// reload haproxy
			cmd := strings.Split(haproxyReloadCmd, " ")
			if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
				return err
			}
		} else if _, err := writeSocket("set weight http-in/" + conf.Servers[0] + " " + strconv.Itoa(256-conf.Subtrahend)); err != nil {
			return err
		} else if _, err := writeSocket("set weight http-in/" + conf.Servers[1] + " " + strconv.Itoa(conf.Subtrahend)); err != nil {
			return err
		}
	}

	return nil
}
