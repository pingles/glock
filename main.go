package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var zookeeper = flag.String("zookeeper", "localhost:2181", "zookeeper connecting string")
var lockPath = flag.String("lockPath", "/some/path", "zookeeper path for lock")

func main() {
	app, err := newApp([]string{*zookeeper}, *lockPath)
	if err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		app.stop()
	}()
	app.run()
}
