package rest

import (
	"bufio"
	"context"
	"errors"
	"strings"

	"github.com/pasarguard/node_bridge/controller"
)

func (n *Node) StreamLogs(ctx context.Context) (<-chan controller.LogEntry, error) {
	if n.Health() == controller.NotConnected {
		return nil, errors.New("node not connected")
	}

	logChan := make(chan controller.LogEntry, n.LogChanSize())

	go func() {
		defer close(logChan)

		// Create a separate client for streaming (timeout = 0)
		client := *n.client
		client.Timeout = 0

		reader, err := n.createStreamingRequest(&client, "GET", "logs")
		if err != nil {
			pushLogEntry(logChan, controller.LogEntry{Err: err})
			return
		}
		defer reader.Close()

		// Close reader when context is done
		go func() {
			select {
			case <-ctx.Done():
			case <-n.ctx.Done():
			}
			reader.Close()
		}()

		bufReader := bufio.NewReader(reader)

		for {
			line, err := bufReader.ReadString('\n')
			if err != nil {
				if ctx.Err() == nil && n.ctx.Err() == nil {
					pushLogEntry(logChan, controller.LogEntry{Err: err})
				}
				return
			}

			line = strings.TrimSpace(line)
			if line != "" {
				pushLogEntry(logChan, controller.LogEntry{Line: line})
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
