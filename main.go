package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
)

var zookeeper = flag.String("zookeeper", "localhost:2181", "zookeeper connecting string")
var lockPath = flag.String("lockPath", "/some/path", "zookeeper path for lock")
var cronSchedule = flag.String("schedule", "* * * * * *", "cron expression for task schedule")
var command = flag.String("command", "/usr/bin/env echo 'hello, world'", "task to execute")

func main() {
	commandAndArgs := *command
	splits := strings.Split(commandAndArgs, " ")
	command := splits[0]
	args := make([]string, 0)
	if len(splits) > 1 {
		args = splits[1:]
	}

	app, err := newApp([]string{*zookeeper}, *lockPath, *cronSchedule, command, args)
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
