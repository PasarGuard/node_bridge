package controller

import (
	"context"
	"sync"
	"time"

	"github.com/pasarguard/node_bridge/common"
)

const (
	MaxChunkSize      = 2000
	DefaultMaxRetries = 5
	InitialBackoff    = 1 * time.Second
	MaxBackoff        = 30 * time.Second
)

type SyncManager struct {
	ctx          context.Context
	syncer       func([]*common.User) error
	pending      map[string]*common.User
	mu           sync.Mutex
	isRunning    bool
	failureCount int
	maxFailures  int
	hardReset    func()
}

func NewSyncManager(ctx context.Context, syncer func([]*common.User) error, hardReset func()) *SyncManager {
	return &SyncManager{
		ctx:         ctx,
		syncer:      syncer,
		pending:     make(map[string]*common.User),
		maxFailures: DefaultMaxRetries,
		hardReset:   hardReset,
	}
}

func (s *SyncManager) UpdateUsers(users []*common.User) {
	s.mu.Lock()
	for _, u := range users {
		s.pending[u.GetEmail()] = u
	}
	if !s.isRunning {
		s.isRunning = true
		go s.Run()
	}
	s.mu.Unlock()
}

func (s *SyncManager) Run() {
	backoff := InitialBackoff

	for {
		s.mu.Lock()
		if len(s.pending) == 0 {
			s.isRunning = false
			s.mu.Unlock()
			return
		}

		// Drain pending users
		users := make([]*common.User, 0, len(s.pending))
		for _, u := range s.pending {
			users = append(users, u)
		}
		// Clear pending map temporarily; we'll requeue failures
		s.pending = make(map[string]*common.User)
		s.mu.Unlock()

		// Process all users (transport handles chunking)
		err := s.syncer(users)

		if err != nil {
			s.failureCount++
			// Requeue failed users (don't overwrite newer updates)
			s.mu.Lock()
			for _, u := range users {
				if _, exists := s.pending[u.GetEmail()]; !exists {
					s.pending[u.GetEmail()] = u
				}
			}
			s.mu.Unlock()

			if s.failureCount >= s.maxFailures {
				if s.hardReset != nil {
					s.hardReset()
				}
				s.failureCount = 0
			}

			// Exponential backoff
			select {
			case <-s.ctx.Done():
				s.mu.Lock()
				s.isRunning = false
				s.mu.Unlock()
				return
			case <-time.After(backoff):
				backoff *= 2
				if backoff > MaxBackoff {
					backoff = MaxBackoff
				}
			}
			continue // Retry with backoff
		}

		// Success
		s.failureCount = 0
		backoff = InitialBackoff
		// Continue loop to check if more users were added during sync
	}
}
