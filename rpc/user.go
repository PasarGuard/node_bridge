package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
)

func (n *Node) SyncUser(baseCtx context.Context) {
mainLoop:
	for {
		select {
		case <-baseCtx.Done():
			return
		default:
		}
		fmt.Println("health:", n.Health())

		switch n.Health() {
		case controller.Broken:
			time.Sleep(5 * time.Second)
			continue mainLoop
		case controller.NotConnected:
			return
		default:
		}

		syncUser, _ := n.client.SyncUser(n.ctx)
	inLoop:
		for {
			select {
			case <-baseCtx.Done():
				return
			case <-syncUser.Context().Done():
				return
			case _, ok := <-n.NotifyChan:
				if !ok {
					return
				}
				continue mainLoop
			case u, ok := <-n.UserChan:
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
	n.mu.Lock()
	defer n.mu.Unlock()

	ctx, cancel := context.WithTimeout(n.ctx, 10*time.Second)
	defer cancel()

	if _, err := n.client.SyncUsers(ctx, &common.Users{Users: users}); err != nil {
		return err
	}
	return nil
}
