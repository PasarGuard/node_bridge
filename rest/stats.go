package rest

import (
	"github.com/m03ed/gozargah_node_bridge/common"
)

func (n *Node) GetSystemStats() (*common.SystemStatsResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.SystemStatsResponse
	err := n.createRequest(n.client, "GET", "stats/system", &common.Empty{}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetBackendStats() (*common.BackendStatsResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.BackendStatsResponse
	err := n.createRequest(n.client, "GET", "stats/backend", &common.Empty{}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetOutboundsStats(reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.StatResponse
	err := n.createRequest(n.client, "GET", "stats/outbounds", &common.StatRequest{Reset_: reset}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetOutboundStats(tag string, reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.StatResponse
	err := n.createRequest(n.client, "GET", "stats/outbound", &common.StatRequest{Name: tag, Reset_: reset}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetInboundsStats(reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.StatResponse
	err := n.createRequest(n.client, "GET", "stats/inbounds", &common.StatRequest{Reset_: reset}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetInboundStats(tag string, reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.StatResponse
	err := n.createRequest(n.client, "GET", "stats/inbound", &common.StatRequest{Name: tag, Reset_: reset}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetUsersStats(reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.StatResponse
	err := n.createRequest(n.client, "GET", "stats/users", &common.StatRequest{Reset_: reset}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetUserStats(email string, reset bool) (*common.StatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.StatResponse
	err := n.createRequest(n.client, "GET", "stats/user", &common.StatRequest{Name: email, Reset_: reset}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetUserOnlineStat(email string) (*common.OnlineStatResponse, error) {
	if err := n.Connected(); err != nil {
		return nil, err
	}

	var stats common.OnlineStatResponse
	err := n.createRequest(n.client, "GET", "stats/user/online", &common.StatRequest{Name: email}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
