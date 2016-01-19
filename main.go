package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var zookeeper = flag.String("zookeeper", "localhost:2181", "zookeeper connecting string")
var lockPath = flag.String("lockPath", "/some/path", "zookeeper path for lock")
var cronSchedule = flag.String("schedule", "* * * * * *", "cron expression for task schedule")
var command = flag.String("command", "echo 'hello, world'", "task to execute")

func main() {
	app, err := newApp([]string{*zookeeper}, *lockPath, *cronSchedule, *command)
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
