package gozargah_node_bridge

import (
	"errors"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
	"github.com/m03ed/gozargah_node_bridge/rpc"
)

type GozargahNode interface {
	Start(string, common.BackendType, []*common.User) error
	Stop()
	NodeVersion() string
	CoreVersion() string
	SyncUsers(users []*common.User) error
	Info() (*common.BaseInfoResponse, error)
	GetSystemStats() (*common.SystemStatsResponse, error)
	GetBackendStats() (*common.BackendStatsResponse, error)
	GetOutboundsStats() (*common.StatResponse, error)
	GetOutboundStats(string) (*common.StatResponse, error)
	GetInboundsStats() (*common.StatResponse, error)
	GetInboundStats(string) (*common.StatResponse, error)
	GetUsersStats() (*common.StatResponse, error)
	GetUserStats(string) (*common.StatResponse, error)
	GetUserOnlineStat(string) (*common.OnlineStatResponse, error)
	GetHealth() controller.Health
	Status() error
	UpdateUser(*common.User) error
	RemoveUser(*common.User) error
	GetLogs() (chan string, error)
}

type NodeProtocol string

const (
	GRPC NodeProtocol = "GRPC"
	REST NodeProtocol = "REST"
)

func NewNode(address string, port int, clientCert, clientKey, serverCA string, nodeProtocol NodeProtocol) (GozargahNode, error) {
	if address == "" {
		return nil, errors.New("address is empty")
	}
	if port == 0 {
		return nil, errors.New("port is empty")
	}

	var node GozargahNode
	var err error
	switch nodeProtocol {
	case GRPC:
		node, err = rpc.NewNode(address, port, clientCert, clientKey, serverCA)
	case REST:
		return nil, errors.New("not implemented yet")
	default:
		return nil, errors.New("unknown node protocol")
	}
	if err != nil {
		return nil, err
	}
	return node, nil
}
