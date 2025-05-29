package controller

import (
	"sync"

	"github.com/google/uuid"

	"github.com/m03ed/gozargah_node_bridge/common"
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

func (c *Controller) Init(apiKey uuid.UUID, extra map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.health = NotConnected
	c.extra = extra
	c.apiKey = apiKey.String()

	c.UserChan = make(chan *common.User)
	c.NotifyChan = make(chan struct{})
	c.LogsChan = make(chan string)
}

func (c *Controller) GetApiKey() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiKey
}

func (c *Controller) GetExtra() map[string]interface{} {
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

func (c *Controller) GetHealth() Health {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.health
}

func (c *Controller) UpdateUser(u *common.User) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.UserChan <- u
	return nil
}

func (c *Controller) GetLogs() (chan string, error) {
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
