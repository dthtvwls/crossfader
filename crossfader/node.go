package crossfader

import (
	"github.com/docker/libkv/store"
)

const defaultNodeConf = `{
  "servers": ["example.com", "example.net"],
  "subtrahend": 0
}`

type Conf struct {
	Servers    [2]string `json:"servers"`
	Subtrahend int       `json:"subtrahend"`
}

type Node struct {
	kv       store.Store
	endpoint string
}

func (node *Node) key() string {
	return "crossfader/" + node.endpoint
}

func (node *Node) get() ([]byte, error) {
	if pair, err := node.kv.Get(node.key()); err != nil {
		return nil, err
	} else {
		return pair.Value, nil
	}
}

func (node *Node) put(conf []byte) error {
	return node.kv.Put(node.key(), conf, nil)
}
