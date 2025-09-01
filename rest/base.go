package rest

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/pasarguard/node_bridge/common"
	"github.com/pasarguard/node_bridge/controller"
	"github.com/pasarguard/node_bridge/tools"
)

type Node struct {
	controller.Controller
	client     *http.Client
	ctx        context.Context
	baseUrl    string
	cancelFunc context.CancelFunc
	mu         sync.Mutex
}

func New(address string, port int, serverCA []byte, apiKey uuid.UUID, logChanSize int, extra map[string]interface{}) (*Node, error) {
	certPool, err := tools.LoadClientPool(serverCA)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	n := &Node{
		Controller: controller.New(apiKey, logChanSize, extra),
		client:     tools.CreateHTTPClient(certPool, address),
		ctx:        ctx,
		baseUrl:    "https://" + net.JoinHostPort(address, fmt.Sprintf("%d", port)),
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

	data := &common.Backend{
		Type:      backendType,
		Config:    config,
		Users:     users,
		KeepAlive: keepAlive,
	}

	n.client.Timeout = time.Second * 15
	var info common.BaseInfoResponse
	if err := n.createRequest(n.client, "POST", "start", data, &info); err != nil {
		return err
	}

	n.Connect(info.GetNodeVersion(), info.GetCoreVersion())
	n.client.Timeout = time.Second * 10

	n.ctx, n.cancelFunc = context.WithCancel(context.Background())

	go n.checkNodeHealth(n.ctx)
	go n.FetchLogs(n.ctx, *n.client)
	go n.SyncUser(n.ctx)

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
	_ = n.createRequest(n.client, "PUT", "stop", &common.Empty{}, &common.Empty{})
}

func (n *Node) Info() (*common.BaseInfoResponse, error) {
	var info common.BaseInfoResponse
	if err := n.createRequest(n.client, "GET", "info", &common.Empty{}, &info); err != nil {
		return nil, err
	}

	return &info, nil
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

func (n *Node) createRequest(client *http.Client, method, endpoint string, data proto.Message, response proto.Message) error {
	body, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, n.baseUrl+"/"+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", n.ApiKey())
	if body != nil {
		req.Header.Set("Content-Type", "application/x-protobuf")
	}

	do, err := client.Do(req)
	if err != nil {
		return err
	}
	defer do.Body.Close()

	responseBody, _ := io.ReadAll(do.Body)
	if err = proto.Unmarshal(responseBody, response); err != nil {
		return err
	}
	return nil
}

func (n *Node) createStreamingRequest(client *http.Client, method, endpoint string) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, n.baseUrl+"/"+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", n.ApiKey())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
