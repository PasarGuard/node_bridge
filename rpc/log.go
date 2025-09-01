package rpc

import (
	"context"
	"time"

	"github.com/pasarguard/node_bridge/common"
	"github.com/pasarguard/node_bridge/controller"
)

func (n *Node) FetchLogs(ctx context.Context) {

mainLoop:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			switch n.Health() {
			case controller.Broken:
				time.Sleep(5 * time.Second)
				continue
			case controller.NotConnected:
				return
			default:
			}

			logsStream, _ := n.client.GetLogs(n.ctx, &common.Empty{})
			for {
				select {
				case <-ctx.Done():
					return
				default:
					logEntry, err := logsStream.Recv()
					if err != nil {
						continue mainLoop
					}
					if logEntry != nil {
						n.LogsChan <- logEntry.GetDetail()
					}
				}
			}
		}
	}
}
