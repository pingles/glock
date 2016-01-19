package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"sync"
	"time"
)

type croney struct {
	lockPath  string
	stopCh    chan bool
	zookeeper []string

	zkConn *zk.Conn
	zkCh   <-chan zk.Event

	currentState zk.State
	sync.Mutex
}

func newApp(zookeeper []string, lockPath string) (*croney, error) {
	stopper := make(chan bool)
	return &croney{
		lockPath:  lockPath,
		zookeeper: zookeeper,
		stopCh:    stopper,
	}, nil
}

func acquireLock(conn *zk.Conn, lockPath string) error {
	log.Println("attempting to acquire lock", lockPath)
	acl := zk.WorldACL(zk.PermAll)
	lock := zk.NewLock(conn, lockPath, acl)
	return lock.Lock()
}

func (c *croney) connectedWithSession() {
	log.Println("connected with session")
	err := acquireLock(c.zkConn, c.lockPath)
	if err == nil {
		log.Println("lock acquired, current process is operational")
	} else {
		log.Println("error creating lock:", err.Error())
	}
}

func (c *croney) handleZkEvent(event zk.Event) {
	switch event.State {
	case zk.StateConnecting:
		log.Println("attempting to connect...")
	case zk.StateConnected:
		log.Println("connected")
	case zk.StateHasSession:
		c.connectedWithSession()
	}
}

func (c *croney) run() {
	conn, ch, err := zk.Connect(c.zookeeper, 5*time.Second)
	if err != nil {
		panic(err)
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

	<-c.stopCh
}

func (c *croney) stop() {
	log.Println("closing connection")
	c.zkConn.Close()
	c.stopCh <- true
}
