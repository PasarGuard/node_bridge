package gozargah_node_bridge

import (
	"errors"

	"github.com/google/uuid"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/controller"
	"github.com/m03ed/gozargah_node_bridge/rest"
	"github.com/m03ed/gozargah_node_bridge/rpc"
)

type GozargahNode interface {
	Start(string, common.BackendType, []*common.User, uint64) error
	Stop()
	NodeVersion() string
	CoreVersion() string
	SyncUsers(users []*common.User) error
	Info() (*common.BaseInfoResponse, error)
	GetSystemStats() (*common.SystemStatsResponse, error)
	GetBackendStats() (*common.BackendStatsResponse, error)
	GetOutboundsStats(bool) (*common.StatResponse, error)
	GetOutboundStats(string, bool) (*common.StatResponse, error)
	GetInboundsStats(bool) (*common.StatResponse, error)
	GetInboundStats(string, bool) (*common.StatResponse, error)
	GetUsersStats(bool) (*common.StatResponse, error)
	GetUserStats(string, bool) (*common.StatResponse, error)
	GetUserOnlineStat(string) (*common.OnlineStatResponse, error)
	GetUserOnlineIpList(string) (*common.StatsOnlineIpListResponse, error)
	GetHealth() controller.Health
	Connected() error
	UpdateUser(*common.User) error
	RemoveUser(*common.User) error
	GetLogs() (chan string, error)
}

type NodeProtocol string

const (
	GRPC NodeProtocol = "GRPC"
	REST NodeProtocol = "REST"
)

func NewNode(address string, port int, serverCA []byte, apiKey uuid.UUID, extra map[string]interface{}, nodeProtocol NodeProtocol) (GozargahNode, error) {
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
		node, err = rpc.NewNode(address, port, serverCA, apiKey, extra)
	case REST:
		node, err = rest.NewNode(address, port, serverCA, apiKey, extra)
	default:
		return nil, errors.New("unknown node protocol")
	}
	if err != nil {
		return nil, err
	}
	return node, nil
}
