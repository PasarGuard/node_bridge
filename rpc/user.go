package rpc

import (
	"context"
	"time"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
)

func (n *Node) SyncUser() {
	baseCtx := n.baseCtx
	notifyChan := n.NotifyChan
	userChan := n.UserChan

mainLoop:
	for {
		select {
		case <-baseCtx.Done():
			return
		default:
		}

		switch n.GetHealth() {
		case controller.Broken:
			time.Sleep(5 * time.Second)
			continue mainLoop
		case controller.NotConnected:
			return
		default:
		}

		syncUser, _ := n.client.SyncUser(n.baseCtx)
	inLoop:
		for {
			select {
			case <-baseCtx.Done():
				return
			case <-syncUser.Context().Done():
				return
			case _, ok := <-notifyChan:
				if !ok {
					return
				}
				continue mainLoop
			case u, ok := <-userChan:
				if !ok {
					return
				}

				if err := syncUser.Send(u); err != nil {
					break inLoop
				}
			}
		}
	}
}

func (n *Node) SyncUsers(users []*common.User) error {
	if err := n.Connected(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 10*time.Second)
	defer cancel()

	if _, err := n.client.SyncUsers(ctx, &common.Users{Users: users}); err != nil {
		return err
	}
	return nil
}
