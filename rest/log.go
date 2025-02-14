package rest

import (
	"bufio"
	"strings"
	"time"

	"github.com/m03ed/gozargah_node_bridge/controller"
)

func (n *Node) FetchLogs() {
	baseCtx := n.baseCtx
	client := *n.client
	client.Timeout = 0
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

			reader, err := n.createStreamingRequest(&client, "GET", "logs")
			if err != nil {
				continue
			}
			defer reader.Close()

			bufReader := bufio.NewReader(reader)

			for {
				line, err := bufReader.ReadString('\n') // Read until newline
				if err != nil {
					_ = reader.Close()
					continue mainLoop
				}

				line = strings.TrimSpace(line)
				if line != "" {
					n.LogsChan <- line
				}
			}
		}
	}
}
