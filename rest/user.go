package rest

import (
	"time"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
)

func (n *Node) SyncUser() {
	baseCtx := n.baseCtx
	for {
		switch n.GetHealth() {
		case controller.Broken:
			time.Sleep(5 * time.Second)
			continue
		case controller.NotConnected:
			return
		default:
		}

		select {
		case <-baseCtx.Done():
			return

		case _, ok := <-n.NotifyChan:
			if !ok {
				return
			}
			continue

		case u, ok := <-n.UserChan:
			if !ok {
				return
			}

			if err := n.createRequest(n.client, "PUT", "/user/sync", u, &common.Empty{}); err != nil {
				continue
			}
		}
	}
}

func (n *Node) SyncUsers(users []*common.User) error {
	if err := n.Connected(); err != nil {
		return err
	}

	data := &common.Users{
		Users: users,
	}

	if err := n.createRequest(n.client, "PUT", "/users/sync", data, &common.Empty{}); err != nil {
		return err
	}

	return nil
}
