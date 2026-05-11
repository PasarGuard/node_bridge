package controller

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/pasarguard/node_bridge/common"
)

func TestSyncManager_Deduplication(t *testing.T) {
	var mu sync.Mutex
	syncedUsers := make(map[string]int)

	syncer := func(users []*common.User) error {
		mu.Lock()
		defer mu.Unlock()
		for _, u := range users {
			syncedUsers[u.GetEmail()]++
		}
		return nil
	}

	sm := NewSyncManager(context.Background(), syncer, nil)

	u1 := &common.User{Email: "user1@example.com"}
	sm.UpdateUsers([]*common.User{u1})
	sm.UpdateUsers([]*common.User{u1})
	sm.UpdateUsers([]*common.User{u1})

	// Wait for sync
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	count := syncedUsers["user1@example.com"]
	mu.Unlock()

	if count != 1 {
		t.Errorf("expected user1 to be synced exactly once due to deduplication, got %d", count)
	}
}

func TestSyncManager_Chunking(t *testing.T) {
	chunkCounts := []int{}
	var mu sync.Mutex

	syncer := func(users []*common.User) error {
		mu.Lock()
		chunkCounts = append(chunkCounts, len(users))
		mu.Unlock()
		return nil
	}

	sm := NewSyncManager(context.Background(), syncer, nil)

	// Update many users
	users := make([]*common.User, 2500)
	for i := 0; i < 2500; i++ {
		users[i] = &common.User{Email: string(rune(i))}
	}
	sm.UpdateUsers(users)

	// Wait for sync
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	totalSynced := 0
	for _, c := range chunkCounts {
		totalSynced += c
	}

	if totalSynced != 2500 {
		t.Errorf("expected 2500 users to be synced, got %d", totalSynced)
	}
}

func TestSyncManager_BackoffAndHardReset(t *testing.T) {
	failCount := 0
	hardResetCalled := false
	var mu sync.Mutex

	syncer := func(users []*common.User) error {
		mu.Lock()
		failCount++
		mu.Unlock()
		return errors.New("temporary failure")
	}

	hardReset := func() {
		mu.Lock()
		hardResetCalled = true
		mu.Unlock()
	}

	sm := NewSyncManager(context.Background(), syncer, hardReset)
	sm.maxFailures = 3

	sm.UpdateUsers([]*common.User{{Email: "fail@example.com"}})

	// Wait for retries and hard reset
	time.Sleep(5 * time.Second)

	mu.Lock()
	defer mu.Unlock()

	if failCount < 3 {
		t.Errorf("expected at least 3 fail attempts, got %d", failCount)
	}
	if !hardResetCalled {
		t.Error("expected hard reset to be called after 3 failures")
	}
}

func TestSyncManager_UpdateUsers_Batching(t *testing.T) {
	var mu sync.Mutex
	syncCount := 0

	syncer := func(users []*common.User) error {
		mu.Lock()
		syncCount++
		mu.Unlock()
		return nil
	}

	sm := NewSyncManager(context.Background(), syncer, nil)

	users := []*common.User{
		{Email: "u1@ex.com"},
		{Email: "u2@ex.com"},
		{Email: "u3@ex.com"},
	}
	sm.UpdateUsers(users)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if syncCount != 1 {
		t.Errorf("expected 1 sync call for batch update, got %d", syncCount)
	}
}
