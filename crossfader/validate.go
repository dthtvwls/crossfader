package crossfader

import (
	"encoding/json"
	"errors"
	"net"
)

func validate(value []byte) (*Conf, error) {
	var conf Conf
	if err := json.Unmarshal(value, &conf); err != nil {
		return nil, err
	} else if len(conf.Servers) != 2 {
		return nil, errors.New("Must have 2 and only 2 servers")
	} else if conf.Servers[0] == conf.Servers[1] {
		return nil, errors.New("Servers must be unique")
	} else if conf.Subtrahend < 0 || conf.Subtrahend > 256 {
		return nil, errors.New("Subtrahend must be between 0 and 256 inclusive")
	} else if _, err := net.LookupIP(conf.Servers[0]); err != nil {
		return nil, err
	} else if _, err := net.LookupIP(conf.Servers[1]); err != nil {
		return nil, err
	} else {
		return &conf, nil
	}
}
