package controller

import (
	"sync"

	"github.com/google/uuid"

	"github.com/pasarguard/node_bridge/common"
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
	apiKey      string
	extra       map[string]interface{}
	mu          sync.RWMutex
}

func New(apiKey uuid.UUID, logChanSize int, extra map[string]interface{}) Controller {
	return Controller{
		health:     NotConnected,
		apiKey:     apiKey.String(),
		extra:      extra,
		UserChan:   make(chan *common.User),
		NotifyChan: make(chan struct{}, 10), // some extra space to avoid deadlock
		LogsChan:   make(chan string, logChanSize),
	}
}

func (c *Controller) ApiKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiKey
}

func (c *Controller) Extra() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.extra
}

func (c *Controller) SetHealth(health Health) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if health == Broken && c.health != Broken {
		c.NotifyChan <- struct{}{}
	}
	c.health = health
}

func (c *Controller) Health() Health {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.health
}

func (c *Controller) UpdateUser(u *common.User) {
	c.UserChan <- u
}

func (c *Controller) Logs() (chan string, error) {
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

	c.UserChan = make(chan *common.User)
	c.NotifyChan = make(chan struct{})
	c.LogsChan = make(chan string)

	c.nodeVersion = ""
	c.coreVersion = ""
	c.health = NotConnected
}
