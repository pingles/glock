package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type lock struct {
	zklock *zk.Lock
}

func (l *lock) Unlock() error {
	return l.zklock.Unlock()
}

func acquireLock(zookeepers []string, sessionTimeout time.Duration, lockPath string) (*lock, error) {
	conn, _, err := zk.Connect(zookeepers, sessionTimeout)
	if err != nil {
		return nil, err
	}

	zkLock := zk.NewLock(conn, lockPath, zk.WorldACL(zk.PermAll))
	err = zkLock.Lock()
	if err != nil {
		return nil, err
	} else {
		return &lock{zkLock}, nil
	}
}
