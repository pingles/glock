package main

import (
	"flag"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

var zookeeper = flag.String("zookeeper", "localhost:2181", "zookeeper connecting string")

func main() {
	_, ch, err := zk.Connect([]string{*zookeeper}, 5*time.Second)
	if err != nil {
		panic(err)
	}

	go func() {
		for m := range ch {
			if m.State == zk.StateConnecting {
				log.Println("attempting to connect...")
			}
		}
	}()

	c := make(chan bool)
	<-c
}
