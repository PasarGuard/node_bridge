package rpc

import (
	"context"
	"time"

	"github.com/m03ed/gozargah_node_bridge/common"
)

func (n *Node) GetSystemStats() (*common.SystemStatsResponse, error) {
	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetSystemStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetBackendStats() (*common.BackendStatsResponse, error) {
	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetBackendStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetStats(reset bool, name string, statType common.StatType) (*common.StatResponse, error) {
	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetStats(ctx, &common.StatRequest{Reset_: reset, Name: name, Type: statType})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUserOnlineStat(email string) (*common.OnlineStatResponse, error) {
	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUserOnlineStats(ctx, &common.StatRequest{Name: email})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUserOnlineIpList(email string) (*common.StatsOnlineIpListResponse, error) {
	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUserOnlineIpListStats(ctx, &common.StatRequest{Name: email})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
