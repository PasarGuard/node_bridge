package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
	"github.com/m03ed/gozargah_node_bridge/tools"
)

type Node struct {
	controller.Controller
	client     common.NodeServiceClient
	baseCtx    context.Context
	cancelFunc context.CancelFunc
	mu         sync.RWMutex
}

func NewNode(address string, port int, clientCert, clientKey, serverCA []byte, extra map[string]interface{}) (*Node, error) {
	tlsConfig, err := tools.CreateTlsConfig(clientCert, clientKey, serverCA)
	if err != nil {
		return nil, err
	}

	target := net.JoinHostPort(address, fmt.Sprintf("%d", port))

	client, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	n := &Node{
		baseCtx:    ctx,
		client:     common.NewNodeServiceClient(client),
		cancelFunc: cancel,
	}
	n.Init(extra)

	return n, nil
}

func (n *Node) Start(config string, backendType common.BackendType, users []*common.User) error {
	if n.GetHealth() != controller.NotConnected {
		n.Stop()
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	req := &common.Backend{
		Type:   backendType,
		Config: config,
		Users:  users,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	info, err := n.client.Start(ctx, req)
	if err != nil {
		return err
	}

	n.Connect(info.GetNodeVersion(), info.GetCoreVersion())

	md := metadata.Pairs("authorization", "Bearer "+info.GetSessionId())
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	n.baseCtx, n.cancelFunc = context.WithCancel(ctx)

	go n.checkNodeHealth()
	go n.SyncUser()
	go n.FetchLogs()

	return nil
}

func (n *Node) Stop() {
	if n.Connected() != nil {
		return
	}
	n.mu.Lock()
	defer n.mu.Unlock()

	existingMD, ok := metadata.FromOutgoingContext(n.baseCtx)
	if !ok {
		existingMD = metadata.MD{}
	}
	n.cancelFunc()

	n.Disconnect()

	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), existingMD), 5*time.Second)
	defer cancel()

	_, _ = n.client.Stop(ctx, nil)
}

func (n *Node) Info() (*common.BaseInfoResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetBaseInfo(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) checkNodeHealth() {
	baseCtx := n.baseCtx
loop:
	for {
		lastHealth := n.GetHealth()
		select {
		case <-baseCtx.Done():
			break loop
		default:
			_, err := n.GetBackendStats()
			switch {
			case err != nil && lastHealth != controller.Broken:
				n.SetHealth(controller.Broken)
			case err == nil && lastHealth != controller.Healthy:
				n.SetHealth(controller.Healthy)
			}
		}
		time.Sleep(2 * time.Second)
	}
}
