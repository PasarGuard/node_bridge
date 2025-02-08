package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/m03ed/gozargah_node_bridge/common"
)

func (n *Node) SyncUser() {
mainLoop:
	for {
		select {
		case <-n.baseCtx.Done():
			return
		default:
		}
		if err := n.Status(); err != nil {
			time.Sleep(5 * time.Second)
			continue mainLoop
		}
		syncUser, _ := n.client.SyncUser(n.baseCtx)
	inLoop:
		for {
			select {
			case <-n.baseCtx.Done():
				return
			case <-syncUser.Context().Done():
				return
			case <-n.NotifyChan:
				continue mainLoop
			case e := <-n.UserChan:
				u, _ := n.GetUser(e)
				if u == nil {
					continue
				}
				if err := syncUser.Send(u); err != nil {
					if errStatus, ok := status.FromError(err); ok {
						switch errStatus.Code() {
						case codes.DeadlineExceeded:
							break inLoop
						case codes.Canceled:
							break inLoop
						default:
							break inLoop
						}
					}
				}
				n.DeleteUser(e)
			}
		}
	}
}

func (n *Node) SyncUsers(users []*common.User) error {
	if err := n.Status(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 10*time.Second)
	defer cancel()

	if _, err := n.client.SyncUsers(ctx, &common.Users{Users: users}); err != nil {
		return err
	}
	return nil
}
