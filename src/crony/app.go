package main

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	Passive = iota
	Active  = iota
)

type crony struct {
	lockPath  string
	task      *cronyTask
	stopCh    chan bool
	zookeeper []string

	zkConn *zk.Conn
	zkCh   <-chan zk.Event

	state int32

	sync.Mutex
}

type cronyTask struct {
	cronSchedule string
	command      string
	args         []string
	directory    string
}

func newApp(zookeeper []string, lockPath string, task *cronyTask) (*crony, error) {
	stopper := make(chan bool)
	return &crony{
		task:      task,
		lockPath:  lockPath,
		zookeeper: zookeeper,
		stopCh:    stopper,
		state:     Passive,
	}, nil
}

func acquireLock(conn *zk.Conn, lockPath string) error {
	log.Println("attempting to acquire lock", lockPath)
	acl := zk.WorldACL(zk.PermAll)
	lock := zk.NewLock(conn, lockPath, acl)
	return lock.Lock()
}

func (c *crony) active() {
	c.Lock()
	defer c.Unlock()
	c.state = Active
}

func (c *crony) isActive() bool {
	return c.state == Active
}

func (c *crony) passive() {
	c.Lock()
	defer c.Unlock()
	c.state = Passive
}

func (c *crony) connectedWithSession() {
	log.Println("connected with session")
	if !c.isActive() {
		err := acquireLock(c.zkConn, c.lockPath)
		if err != nil {
			log.Println("error acquiring lock:", err.Error())
		}

		if err == nil {
			log.Println("lock acquired, current process is operational")
			c.active()
		}
	}
}

func (c *crony) handleZkEvent(event zk.Event) {
	switch event.State {
	case zk.StateDisconnected:
		// disconnected
	case zk.StateConnecting:
		log.Println("attempting to connect...")
	case zk.StateConnected:
		log.Println("connected")
	case zk.StateHasSession:
		c.connectedWithSession()
	}
}

func runCommand(task *cronyTask) error {
	cmd := exec.Command(task.command, task.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if task.directory != "" {
		cmd.Dir = task.directory
	}
	err := cmd.Run()
	if err != nil {
		log.Println("error running command:", err.Error())
		return err
	}

	return nil
}

func (c *crony) executeTask() {
	if c.isActive() {
		log.Println(fmt.Sprintf("executing command=%s, args=%s, dir=%s", c.task.command, c.task.args, c.task.directory))
		runCommand(c.task)
	} else {
		log.Println("passive, ignoring command execution", c.task.command)
	}
}

func (c *crony) run() {
	sessionTimeout := 5 * time.Second
	conn, ch, err := zk.Connect(c.zookeeper, sessionTimeout)
	if err != nil {
		log.Fatal(err)
	}
	c.zkConn = conn
	c.zkCh = ch

	log.Println("starting")
	go func() {
		for {
			m := <-c.zkCh
			c.handleZkEvent(m)
		}
	}()

	cron := cron.New()
	err = cron.AddFunc(c.task.cronSchedule, c.executeTask)
	if err != nil {
		log.Fatal(err)
	}
	cron.Start()

	<-c.stopCh
}

func (c *crony) stop() {
	log.Println("closing connection")
	c.zkConn.Close()
	c.stopCh <- true
}
