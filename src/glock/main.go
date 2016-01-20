package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"time"
)

var (
	zookeeper      = kingpin.Flag("zookeeper", "zookeeper connection string. can be comma-separated.").Required().String()
	path           = kingpin.Flag("path", "zookeeper path for lock").Required().String()
	command        = kingpin.Flag("command", "command to execute").Required().String()
	sleep          = kingpin.Flag("sleep", "duration to sleep after task executes to ensure single execution.").Default("15s").Duration()
	sessionTimeout = kingpin.Flag("sessionTimeout", "zookeeper session timeout.").Default("5s").Duration()
)

func parseZooKeeper(zkarg string) []string {
	return strings.Split(zkarg, ",")
}

func main() {
	kingpin.Parse()

	zookeepers := parseZooKeeper(*zookeeper)
	lock, err := acquireLock(zookeepers, *sessionTimeout, *path)
	if err != nil {
		fmt.Println("error acquiring lock, exiting:", err.Error())
		os.Exit(1)
	}

	if lock != nil {
		fmt.Println("LOCK ACQUIRED, EXECUTING")
		runCommand(parseCommand(*command))
	} else {
		fmt.Println("error acquiring lock, exiting")
		os.Exit(1)
	}

	fmt.Println("waiting for", *sleep)
	<-time.After(*sleep)

	lock.Unlock()
}
