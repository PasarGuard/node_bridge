package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
	"github.com/m03ed/gozargah_node_bridge/tools"
)

type Node struct {
	controller.Controller
	client     *http.Client
	baseCtx    context.Context
	baseUrl    string
	sessionId  string
	cancelFunc context.CancelFunc
	mu         sync.RWMutex
}

func NewNode(address string, port int, clientCert, clientKey, serverCA []byte) (*Node, error) {
	tlsConfig, err := tools.CreateTlsConfig(clientCert, clientKey, serverCA)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	n := &Node{
		client:     tools.CreateHTTPClient(tlsConfig),
		baseCtx:    ctx,
		baseUrl:    "https://" + net.JoinHostPort(address, fmt.Sprintf("%d", port)),
		sessionId:  "",
		cancelFunc: cancel,
	}
	n.Init()

	return n, nil
}

func (n *Node) Start(config string, backendType common.BackendType, users []*common.User) error {
	if n.GetHealth() != controller.NotConnected {
		n.Stop()
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	data := &common.Backend{
		Type:   backendType,
		Config: config,
		Users:  users,
	}

	n.client.Timeout = time.Second * 15
	var info common.BaseInfoResponse
	if err := n.createRequest(n.client, "POST", "/start", data, &info); err != nil {
		return err
	}

	n.Connect(info.GetNodeVersion(), info.GetCoreVersion())
	n.sessionId = info.GetSessionId()
	n.client.Timeout = time.Second * 5

	n.baseCtx, n.cancelFunc = context.WithCancel(context.Background())

	go n.checkNodeHealth()
	go n.FetchLogs()
	go n.SyncUser()

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
	_ = n.createRequest(n.client, "PUT", "/stop", &common.Empty{}, &common.Empty{})

	n.sessionId = ""
}

func (n *Node) Info() (*common.BaseInfoResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var info common.BaseInfoResponse
	if err := n.createRequest(n.client, "GET", "/info", &common.Empty{}, &info); err != nil {
		return nil, err
	}

	return &info, nil
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

func (n *Node) createRequest(client *http.Client, method, endpoint string, data proto.Message, response proto.Message) error {
	body, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, n.baseUrl+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+n.sessionId)
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
