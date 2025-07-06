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
	GetStats(reset bool, name string, statType common.StatType) (*common.StatResponse, error)
	GetUserOnlineStat(string) (*common.OnlineStatResponse, error)
	GetUserOnlineIpList(string) (*common.StatsOnlineIpListResponse, error)
	Health() controller.Health
	UpdateUser(*common.User)
	Logs() (chan string, error)
}

type NodeProtocol string

const (
	GRPC NodeProtocol = "GRPC"
	REST NodeProtocol = "REST"
)

// NodeOptions holds the configuration for creating a new node
type NodeOptions struct {
	address      string
	port         int
	serverCA     []byte
	apiKey       uuid.UUID
	extra        map[string]interface{}
	nodeProtocol NodeProtocol
	logChanSize  int
}

// NodeOption is a function type for configuring NodeOptions
type NodeOption func(*NodeOptions) error

// WithPort sets the port for the node
func WithPort(port int) NodeOption {
	return func(opts *NodeOptions) error {
		if port <= 0 {
			return errors.New("port must be greater than 0")
		}
		opts.port = port
		return nil
	}
}

// WithServerCA sets the server CA certificate
func WithServerCA(serverCA []byte) NodeOption {
	return func(opts *NodeOptions) error {
		opts.serverCA = serverCA
		return nil
	}
}

// WithAPIKey sets the API key
func WithAPIKey(apiKey uuid.UUID) NodeOption {
	return func(opts *NodeOptions) error {
		opts.apiKey = apiKey
		return nil
	}
}

// WithExtra sets extra configuration parameters
func WithExtra(extra map[string]interface{}) NodeOption {
	return func(opts *NodeOptions) error {
		opts.extra = extra
		return nil
	}
}

func WithLogChannelSize(size int) NodeOption {
	return func(opts *NodeOptions) error {
		if size <= 0 {
			opts.logChanSize = 1000
		} else {
			opts.logChanSize = size
		}
		return nil
	}
}

// New creates a new node with the given address, protocol, and options
func New(address string, nodeProtocol NodeProtocol, options ...NodeOption) (GozargahNode, error) {
	if address == "" {
		return nil, errors.New("address is empty")
	}

	// Initialize options with defaults
	opts := &NodeOptions{
		address:      address,
		nodeProtocol: nodeProtocol,
		extra:        make(map[string]interface{}),
	}

	// Apply all provided options
	for _, option := range options {
		if err := option(opts); err != nil {
			return nil, err
		}
	}

	var node GozargahNode
	var err error
	switch nodeProtocol {
	case GRPC:
		node, err = rpc.New(opts.address, opts.port, opts.serverCA, opts.apiKey, opts.logChanSize, opts.extra)
	case REST:
		node, err = rest.New(opts.address, opts.port, opts.serverCA, opts.apiKey, opts.logChanSize, opts.extra)
	default:
		return nil, errors.New("unknown node protocol")
	}
	if err != nil {
		return nil, err
	}
	return node, nil
}
