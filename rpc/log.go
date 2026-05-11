package rpc

import (
	"context"
	"errors"

	"github.com/pasarguard/node_bridge/common"
	"github.com/pasarguard/node_bridge/controller"
)

func (n *Node) StreamLogs(ctx context.Context) (<-chan controller.LogEntry, error) {
	if n.Health() == controller.NotConnected {
		return nil, errors.New("node not connected")
	}

	logChan := make(chan controller.LogEntry, n.LogChanSize())

	go func() {
		defer close(logChan)

		logsStream, err := n.client.GetLogs(n.ctx, &common.Empty{})
		if err != nil {
			pushLogEntry(logChan, controller.LogEntry{Err: err})
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-n.ctx.Done():
				return
			default:
				logEntry, err := logsStream.Recv()
				if err != nil {
					// Only push error if it's not a normal cancellation
					if ctx.Err() == nil && n.ctx.Err() == nil {
						pushLogEntry(logChan, controller.LogEntry{Err: err})
					}
					return
				}
				if logEntry != nil {
					pushLogEntry(logChan, controller.LogEntry{Line: logEntry.GetDetail()})
				}
			}
		}
	}()

	return logChan, nil
}

func pushLogEntry(ch chan controller.LogEntry, entry controller.LogEntry) {
	select {
	case ch <- entry:
	default:
		// Drop oldest
		select {
		case <-ch:
		default:
		}
		// Non-blocking write to avoid deadlock if channel was emptied between selects
		select {
		case ch <- entry:
		default:
		}
	}
}
