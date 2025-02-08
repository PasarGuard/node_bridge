package rpc

import (
	"context"
	"time"

	"github.com/m03ed/gozargah_node_bridge/common"
)

func (n *Node) GetSystemStats() (*common.SystemStatsResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetSystemStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetBackendStats() (*common.BackendStatsResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetBackendStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetOutboundsStats() (*common.StatResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetOutboundsStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetOutboundStats(tag string) (*common.StatResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetOutboundStats(ctx, &common.StatRequest{Name: tag})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetInboundsStats() (*common.StatResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetInboundsStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetInboundStats(tag string) (*common.StatResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetInboundStats(ctx, &common.StatRequest{Name: tag})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUsersStats() (*common.StatResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUsersStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUserStats(email string) (*common.StatResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUserStats(ctx, &common.StatRequest{Name: email})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUserOnlineStat(email string) (*common.OnlineStatResponse, error) {
	if err := n.Status(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUserOnlineStats(ctx, &common.StatRequest{Name: email})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
