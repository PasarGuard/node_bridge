# PasarGuard Node Bridge

A Go library for To connect easy and stable to [PasarGuard-Node](https://github.com/PasarGuard/node).

## Installation
```bash
go get github.com/pasarguard/node_bridge
```

## Usage Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/pasarguard/node_bridge"
	"github.com/pasarguard/node_bridge/common"
)

func main() {
	// Initialize the node client
	apiKey := uuid.New() // Replace with your actual API key
	node, err := node_bridge.New("127.0.0.1", node_bridge.GRPC,
		node_bridge.WithPort(8080),
		node_bridge.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Start the node with a configuration
	config := `{"inbounds": [], "outbounds": []}` // Your node configuration
	err = node.Start(config, common.BackendType_XRAY, nil, 60)
	if err != nil {
		log.Fatal(err)
	}

	// Get node information
	info, err := node.Info()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Node Version: %s\n", info.NodeVersion)
	fmt.Printf("Core Version: %s\n", info.CoreVersion)

	// Get system stats
	stats, err := node.GetSystemStats()
	if err == nil {
		fmt.Printf("CPU Usage: %.2f%%\n", stats.CpuUsage)
	}
}
```
