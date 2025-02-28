package rpc

import (
	"context"
	"time"

	"github.com/m03ed/gozargah_node_bridge/common"
)

func (n *Node) GetSystemStats() (*common.SystemStatsResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetSystemStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetBackendStats() (*common.BackendStatsResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetBackendStats(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetOutboundsStats(reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetOutboundsStats(ctx, &common.StatRequest{Reset_: reset})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetOutboundStats(tag string, reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetOutboundStats(ctx, &common.StatRequest{Name: tag, Reset_: reset})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetInboundsStats(reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetInboundsStats(ctx, &common.StatRequest{Reset_: reset})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetInboundStats(tag string, reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetInboundStats(ctx, &common.StatRequest{Name: tag, Reset_: reset})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUsersStats(reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUsersStats(ctx, &common.StatRequest{Reset_: reset})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUserStats(email string, reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUserStats(ctx, &common.StatRequest{Name: email, Reset_: reset})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (n *Node) GetUserOnlineStat(email string) (*common.OnlineStatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, cancel := context.WithTimeout(n.baseCtx, 5*time.Second)
	defer cancel()

	resp, err := n.client.GetUserOnlineStats(ctx, &common.StatRequest{Name: email})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
