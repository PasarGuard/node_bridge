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

	"github.com/pasarguard/node_bridge/common"
	"github.com/pasarguard/node_bridge/controller"
	"github.com/pasarguard/node_bridge/tools"
)

type Node struct {
	controller.Controller
	client     common.NodeServiceClient
	ctx        context.Context
	cancelFunc context.CancelFunc
	mu         sync.Mutex
}

func New(address string, port int, serverCA []byte, apiKey uuid.UUID, logChanSize int, extra map[string]interface{}) (*Node, error) {
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

	ctx, cancel := createCtxWithMD(apiKey.String())

	n := &Node{
		Controller: controller.New(apiKey, logChanSize, extra),
		ctx:        ctx,
		client:     common.NewNodeServiceClient(client),
		cancelFunc: cancel,
	}

	return n, nil
}

func (n *Node) Start(config string, backendType common.BackendType, users []*common.User, keepAlive uint64) error {
	if n.Health() != controller.NotConnected {
		n.Stop()
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	req := &common.Backend{
		Type:      backendType,
		Config:    config,
		Users:     users,
		KeepAlive: keepAlive,
	}

	ctx, cancel := context.WithTimeout(n.ctx, 15*time.Second)
	defer cancel()

	info, err := n.client.Start(ctx, req)
	if err != nil {
		return err
	}

	n.Connect(info.GetNodeVersion(), info.GetCoreVersion())

	go n.checkNodeHealth(n.ctx)
	go n.SyncUser(n.ctx)
	go n.FetchLogs(n.ctx)

	return nil
}

func (n *Node) Stop() {
	if n.Health() == controller.NotConnected {
		return
	}
	n.mu.Lock()
	defer n.mu.Unlock()

	n.cancelFunc()
	n.Disconnect()

	n.ctx, n.cancelFunc = createCtxWithMD(n.ApiKey())

	ctx, cancel := context.WithTimeout(n.ctx, 5*time.Second)
	defer cancel()

	_, _ = n.client.Stop(ctx, nil)
}

func (n *Node) Info() (*common.BaseInfoResponse, error) {
	ctx, cancel := context.WithTimeout(n.ctx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetBaseInfo(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) checkNodeHealth(ctx context.Context) {
loop:
	for {
		lastHealth := n.Health()
		select {
		case <-ctx.Done():
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

func createCtxWithMD(apiKey string) (context.Context, context.CancelFunc) {
	md := metadata.Pairs("x-api-key", apiKey)
	ctxWithKey := metadata.NewOutgoingContext(context.Background(), md)
	return context.WithCancel(ctxWithKey)
}
