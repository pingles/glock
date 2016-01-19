package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"github.com/robfig/cron"
	"log"
	"sync"
	"time"
)

const (
	Passive = iota
	Active = iota
)

type croney struct {
	cronSchedule string
	command      string
	lockPath     string
	stopCh       chan bool
	zookeeper    []string

	zkConn *zk.Conn
	zkCh   <-chan zk.Event

	state int32

	sync.Mutex
}

func newApp(zookeeper []string, lockPath, schedule, command string) (*croney, error) {
	stopper := make(chan bool)
	return &croney{
		cronSchedule: schedule,
		command:      command,
		lockPath:     lockPath,
		zookeeper:    zookeeper,
		stopCh:       stopper,
		state:        Passive,
	}, nil
}

func acquireLock(conn *zk.Conn, lockPath string) error {
	log.Println("attempting to acquire lock", lockPath)
	acl := zk.WorldACL(zk.PermAll)
	lock := zk.NewLock(conn, lockPath, acl)
	return lock.Lock()
}

func (c *croney) active() {
	c.Lock()
	defer c.Unlock()
	c.state = Active
}

func (c *croney) isActive() bool {
	return c.state == Active
}

func (c *croney) passive() {
	c.Lock()
	defer c.Unlock()
	c.state = Passive
}

func (c *croney) connectedWithSession() {
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

func (c *croney) handleZkEvent(event zk.Event) {
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

func (c *croney) executeTask() {
	if c.isActive() {
		log.Println("executing task:", c.command)
	} else {
		log.Println("not active, won't run command", c.command)
	}
}

func (c *croney) run() {
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
	err = cron.AddFunc(c.cronSchedule, c.executeTask)
	if err != nil {
		log.Fatal(err)
	}
	cron.Start()

	<-c.stopCh
}

func (c *croney) stop() {
	log.Println("closing connection")
	c.zkConn.Close()
	c.stopCh <- true
}
