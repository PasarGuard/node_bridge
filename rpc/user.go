package rpc

import (
	"context"
	"time"

	"github.com/pasarguard/node_bridge/common"
	"github.com/pasarguard/node_bridge/controller"
)

func (n *Node) SyncUsers(users []*common.User) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	ctx, cancel := context.WithTimeout(n.ctx, 30*time.Second)
	defer cancel()

	stream, err := n.client.SyncUsersChunked(ctx)
	if err != nil {
		return err
	}

	for i := 0; i < len(users); i += controller.MaxChunkSize {
		end := i + controller.MaxChunkSize
		if end > len(users) {
			end = len(users)
		}
		chunk := users[i:end]

		err := stream.Send(&common.UsersChunk{
			Users: chunk,
			Index: uint64(i / controller.MaxChunkSize),
			Last:  end == len(users),
		})
		if err != nil {
			return err
		}
	}

	if len(users) == 0 {
		err := stream.Send(&common.UsersChunk{
			Users: nil,
			Index: 0,
			Last:  true,
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	return err
}
