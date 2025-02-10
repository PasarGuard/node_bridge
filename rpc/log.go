package rpc

import (
	"errors"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
)

func (n *Node) FetchLogs() {
	baseCtx := n.baseCtx

mainLoop:
	for {
		select {
		case <-baseCtx.Done():
			return
		default:
			switch n.GetHealth() {
			case controller.Broken:
				time.Sleep(5 * time.Second)
				continue
			case controller.NotConnected:
				return
			default:
			}

			logsStream, _ := n.client.GetLogs(n.baseCtx, &common.Empty{})
			for {
				select {
				case <-baseCtx.Done():
					return
				default:
					logEntry, err := logsStream.Recv()
					if err != nil {
						_ = logsStream.CloseSend()
						if errors.Is(err, io.EOF) {
							continue mainLoop
						}

						if st, ok := status.FromError(err); ok {
							switch st.Code() {
							case codes.Canceled, codes.Unavailable, codes.DeadlineExceeded:
								continue mainLoop
							default:
								n.SetHealth(controller.Broken)
							}
						}
					}
					if logEntry != nil {
						n.LogsChan <- logEntry.GetDetail()
					}
				}
			}
		}
	}
}
