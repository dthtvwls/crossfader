package crossfader

import (
	"io/ioutil"
  "log"
	"net/http"
  "reflect"
	"strings"

	"github.com/docker/libkv"
  "github.com/docker/libkv/store"
)

func Start(backend, server, endpoint string) error {
	if kv, err := libkv.NewStore(store.Backend(backend), strings.Split(server, ","), &store.Config{}); err != nil {
		return err
	} else {
		node := Node{kv: kv, endpoint: endpoint}

    // ensure key exists
		if exists, err := kv.Exists(node.key()); err != nil {
			return err
		} else if !exists {
			if err := node.put([]byte(defaultNodeConf)); err != nil {
				return err
      }
		}

    // watch key for updates
		go func() {
    	if events, err := kv.Watch(node.key(), make(<-chan struct{})); err != nil {
    		log.Fatal(err)
    	} else {
    		for {
    			pair := <-events

    			if err := putHaproxyConf(pair.Value); err != nil {
    				log.Print(err)
    			}
    		}
    	}
    }()

    http.HandleFunc("/crossfader", func(w http.ResponseWriter, r *http.Request) {
      if r.Method == "PUT" {
        defer r.Body.Close()
        if conf, err := ioutil.ReadAll(r.Body); err != nil {
          http.Error(w, err.Error(), http.StatusBadRequest)
        } else if _, err := validate(conf); err != nil {
          http.Error(w, err.Error(), http.StatusBadRequest)
        } else if err := node.put(conf); err != nil {
          http.Error(w, err.Error(), http.StatusBadGateway)
        } else {
          w.Write(conf)
        }
      } else if conf, err := node.get(); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
      } else {
        w.Write(conf)
      }
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      if value, err := node.get(); err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
      } else if conf, err := validate(value); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
      } else if haproxyConf, err := getHaproxyConf(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
      } else if !reflect.DeepEqual(conf, haproxyConf) {
        http.Error(w, "HAProxy misconfigured", http.StatusInternalServerError)
      } else {
        w.Write([]byte("OK"))
      }
    })

    http.Handle("/", http.FileServer(http.Dir("public")))

		return http.ListenAndServe(":2468", nil)
	}
}
