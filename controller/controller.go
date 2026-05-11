package controller

import (
	"context"
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

type LogEntry struct {
	Line string
	Err  error
}

type Controller struct {
	health        Health
	nodeVersion   string
	coreVersion   string
	apiKey        string
	extra         map[string]interface{}
	logChanSize   int
	mu            sync.RWMutex
	SyncManager   *SyncManager
	HardResetChan chan struct{}
}

func New(apiKey uuid.UUID, logChanSize int, extra map[string]interface{}) Controller {
	return Controller{
		health:        NotConnected,
		apiKey:        apiKey.String(),
		extra:         extra,
		logChanSize:   logChanSize,
		HardResetChan: make(chan struct{}, 1),
	}
}

func (c *Controller) LogChanSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logChanSize
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
	c.health = health
}

func (c *Controller) Health() Health {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.health
}

func (c *Controller) UpdateUsers(users []*common.User) {
	c.mu.RLock()
	sm := c.SyncManager
	c.mu.RUnlock()
	if sm != nil {
		sm.UpdateUsers(users)
	}
}

func (c *Controller) HardReset() <-chan struct{} {
	return c.HardResetChan
}

func (c *Controller) triggerHardReset() {
	select {
	case c.HardResetChan <- struct{}{}:
	default:
	}
}

func (c *Controller) StartSync(ctx context.Context, syncer func([]*common.User) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SyncManager = NewSyncManager(ctx, syncer, c.triggerHardReset)
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

	close(c.HardResetChan)

	c.HardResetChan = make(chan struct{}, 1)
	c.SyncManager = nil

	c.nodeVersion = ""
	c.coreVersion = ""
	c.health = NotConnected
}
