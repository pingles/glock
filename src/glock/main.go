package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	zookeeper      = kingpin.Flag("zookeeper", "zookeeper connection string. can be comma-separated.").Required().String()
	path           = kingpin.Flag("path", "zookeeper path for lock").Required().String()
	command        = kingpin.Flag("command", "command to execute").Required().String()
	minExec        = kingpin.Flag("minExec", "minimum execution time. should be set large enough to cover clock drift across instances.").Default("15s").Duration()
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

func execChannel(cmd *commandAndArgs) <-chan error {
	c := make(chan error, 1)
	go func() {
		err := runCommand(cmd)
		c <- err
	}()
	return c
}

func main() {
	kingpin.Parse()

	lockCh := lockChannel(*zookeeper, *sessionTimeout, *path)

	select {
	case lock := <-lockCh:
		execCh := execChannel(parseCommand(*command))
		<-time.After(*minExec)
		execError := <-execCh

		lock.Unlock()

		if execError != nil {
			fmt.Println("error executing command:", execError.Error())
			exitError := execError.(*exec.ExitError)
			if exitError != nil {
				exitCode := exitError.Sys().(syscall.WaitStatus).ExitStatus()
				os.Exit(exitCode)
			}
		}
	case <-time.After(*wait):
		fmt.Println("couldn't acquire lock in time, exiting.")
		os.Exit(1)
	}
}
