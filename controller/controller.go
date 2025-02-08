package controller

import (
	"errors"
	"sync"

	"github.com/m03ed/gozargah_node_bridge/common"
)

var (
	NotConnectedError = errors.New("node is not connected")
	BrokenNodeError   = errors.New("node is broken")
	UserNotFoundError = errors.New("user not found")
	DisconnectedError = errors.New("node is disconnected")
)

type Health int

const (
	NotConnected Health = iota
	Broken
	Healthy
	Disconnected
)

type Controller struct {
	health     Health
	users      map[string]*common.User
	UserChan   chan string
	NotifyChan chan struct{}
	LogsChan   chan string
	mu         sync.RWMutex
}

func (c *Controller) Init() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.users = make(map[string]*common.User)
	c.UserChan = make(chan string)
	c.NotifyChan = make(chan struct{})
	c.LogsChan = make(chan string)
	c.health = NotConnected
}

func (c *Controller) SetHealth(health Health) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.health != Broken && health == Broken {
		c.NotifyChan <- struct{}{}
	}
	c.health = health
}

func (c *Controller) GetHealth() Health {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.health
}

func (c *Controller) Status() error {
	switch c.GetHealth() {
	case NotConnected:
		return NotConnectedError
	case Broken:
		return BrokenNodeError
	case Disconnected:
		return DisconnectedError
	default:
		return nil
	}
}

func (c *Controller) UpdateUser(u *common.User) error {
	if err := c.Status(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.users[u.GetEmail()] = u
	c.UserChan <- u.GetEmail()
	return nil
}

func (c *Controller) RemoveUser(u *common.User) error {
	if err := c.Status(); err != nil {
		return err
	}
	u.Inbounds = []string{}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.users[u.GetEmail()] = u
	c.UserChan <- u.GetEmail()

	return nil
}

func (c *Controller) GetUser(e string) (*common.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	user, ok := c.users[e]
	if !ok {
		return nil, UserNotFoundError
	}
	return user, nil
}

func (c *Controller) DeleteUser(e string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.users, e)
}

func (c *Controller) GetLogs() (chan string, error) {
	switch c.GetHealth() {
	case NotConnected:
		return nil, NotConnectedError
	case Disconnected:
		return nil, DisconnectedError
	default:
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.LogsChan, nil
}
