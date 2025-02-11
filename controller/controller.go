package controller

import (
	"errors"
	"sync"

	"github.com/m03ed/gozargah_node_bridge/common"
)

var (
	NotConnectedError = errors.New("node is not connected")
)

type Health int

const (
	NotConnected Health = iota
	Broken
	Healthy
)

type Controller struct {
	health      Health
	UserChan    chan *common.User
	NotifyChan  chan struct{}
	LogsChan    chan string
	nodeVersion string
	coreVersion string
	mu          sync.RWMutex
}

func (c *Controller) Init() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.health = NotConnected
}

func (c *Controller) SetHealth(health Health) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if health == Broken && c.health != Broken {
		c.NotifyChan <- struct{}{}
	}
	c.health = health
}

func (c *Controller) GetHealth() Health {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.health
}

func (c *Controller) Connected() error {
	switch c.GetHealth() {
	case NotConnected:
		return NotConnectedError
	default:
		return nil
	}
}

func (c *Controller) UpdateUser(u *common.User) error {
	if err := c.Connected(); err != nil {
		return err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	c.UserChan <- u
	return nil
}

func (c *Controller) RemoveUser(u *common.User) error {
	if err := c.Connected(); err != nil {
		return err
	}
	u.Inbounds = []string{}

	c.mu.RLock()
	defer c.mu.RUnlock()

	c.UserChan <- u
	return nil
}

func (c *Controller) GetLogs() (chan string, error) {
	switch c.GetHealth() {
	case NotConnected:
		return nil, NotConnectedError
	default:
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.LogsChan, nil
}

func (c *Controller) NodeVersion() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nodeVersion
}

func (c *Controller) CoreVersion() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.coreVersion
}

func (c *Controller) Connect(nodeVersion, coreVersion string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.UserChan = make(chan *common.User)
	c.NotifyChan = make(chan struct{})
	c.LogsChan = make(chan string)

	c.nodeVersion = nodeVersion
	c.coreVersion = coreVersion
	c.health = Healthy
}

func (c *Controller) Disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.UserChan)
	close(c.NotifyChan)
	close(c.LogsChan)

	c.UserChan = nil
	c.NotifyChan = nil
	c.LogsChan = nil

	c.nodeVersion = ""
	c.coreVersion = ""
	c.health = NotConnected
}
