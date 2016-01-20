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
	wait           = kingpin.Flag("wait", "duration to wait for lock before exiting.").Default("5s").Duration()
	sessionTimeout = kingpin.Flag("sessionTimeout", "zookeeper session timeout.").Default("10s").Duration()
)

func parseZooKeeper(zkarg string) []string {
	return strings.Split(zkarg, ",")
}

func lockChannel(zookeeper string, sessionTimeout time.Duration, path string) <-chan *lock {
	lockCh := make(chan *lock, 1)
	go func() {
		zookeepers := parseZooKeeper(zookeeper)
		lock, err := acquireLock(zookeepers, sessionTimeout, path)
		if err != nil {
			fmt.Println("error acquiring lock, exiting:", err.Error())
			os.Exit(1)
		}

		if lock == nil {
			fmt.Println("didn't acquire lock")
			os.Exit(1)
		}

		if lock != nil {
			lockCh <- lock
		}
	}()
	return lockCh
}

func main() {
	kingpin.Parse()

	lockCh := lockChannel(*zookeeper, *sessionTimeout, *path)

	select {
	case lock := <-lockCh:
		runCommand(parseCommand(*command))
		time.Sleep(*sleep)
		lock.Unlock()
	case <-time.After(*wait):
		fmt.Println("couldn't acquire lock in time, exiting.")
		os.Exit(1)
	}
}
