package rest

import (
	"github.com/m03ed/gozargah_node_bridge/controller"
	"time"

	"github.com/m03ed/gozargah_node_bridge/common"
)

func (n *Node) FetchLogs() {
	baseCtx := n.baseCtx
	client := *n.client
	client.Timeout = 60 * time.Second
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

			var logs common.LogList
			if err := n.createRequest(&client, "GET", "/logs", &common.Empty{}, &logs); err != nil {
				continue
			}

			for _, logEntry := range logs.GetLogs() {
				if logEntry != "" {
					n.LogsChan <- logEntry
				}
			}
		}
	}
}
