package rpc

import (
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
)

func (n *Node) FetchLogs() {
mainLoop:
	for {
		select {
		case <-n.baseCtx.Done():
			return
		default:
			if err := n.Status(); err != nil {
				time.Sleep(5 * time.Second)
				continue mainLoop
			}
			logsStream, _ := n.client.GetLogs(n.baseCtx, &common.Empty{})
			for {
				select {
				case <-n.baseCtx.Done():
					return
				default:
					logEntry, err := logsStream.Recv()
					if err != nil {
						_ = logsStream.CloseSend()
						time.Sleep(5 * time.Second)
						if err == io.EOF {
							continue mainLoop
						}

						if st, ok := status.FromError(err); ok {
							switch st.Code() {
							case codes.Canceled:
								continue mainLoop
							case codes.Unavailable, codes.DeadlineExceeded:
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
