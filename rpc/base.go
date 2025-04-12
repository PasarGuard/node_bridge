package rpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

func NewNode(address string, port int, serverCA []byte, apiKey uuid.UUID, extra map[string]interface{}) (*Node, error) {
	certPool, err := tools.LoadClientPool(serverCA)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewClientTLSFromCert(certPool, "")
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	target := net.JoinHostPort(address, fmt.Sprintf("%d", port))

	client, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %v", err)
	}

	md := metadata.Pairs("authorization", "Bearer "+apiKey.String())
	ctxWithKey := metadata.NewOutgoingContext(context.Background(), md)
	ctx, cancel := context.WithCancel(ctxWithKey)

	n := &Node{
		baseCtx:    ctx,
		client:     common.NewNodeServiceClient(client),
		cancelFunc: cancel,
	}
	n.Init(apiKey, extra)

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

	ctx, cancel := context.WithTimeout(n.baseCtx, 15*time.Second)
	defer cancel()

	info, err := n.client.Start(ctx, req)
	if err != nil {
		return err
	}

	n.Connect(info.GetNodeVersion(), info.GetCoreVersion())

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

	n.cancelFunc()
	n.Disconnect()

	md := metadata.Pairs("authorization", "Bearer "+n.GetApiKey())
	ctxWithKey := metadata.NewOutgoingContext(context.Background(), md)

	ctx, cancel := context.WithTimeout(ctxWithKey, 5*time.Second)
	defer cancel()

	_, _ = n.client.Stop(ctx, nil)

	n.baseCtx, n.cancelFunc = context.WithCancel(ctxWithKey)
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
