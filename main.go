package main

import (
	"flag"
	"log"
	"os"

	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/dthtvwls/crossfader/crossfader"
)

var backend, server, endpoint string

func init() {
	flag.StringVar(&backend, "backend", os.Getenv("BACKEND"), "Backend type (consul/etcd/zk)")
	flag.StringVar(&server, "server", os.Getenv("SERVER"), "Backend server(s), comma-separated")
	flag.StringVar(&endpoint, "endpoint", os.Getenv("ENDPOINT"), "DNS endpoint")
	flag.Parse()

	switch backend {
	case "consul":
		consul.Register()
	case "etcd":
		etcd.Register()
	case "zk":
		zookeeper.Register()
	}
}

func main() {
	log.Fatal(crossfader.Start(backend, server, endpoint))
}
