package rest

import (
	"github.com/pasarguard/node_bridge/common"
)

func (n *Node) GetSystemStats() (*common.SystemStatsResponse, error) {
	var stats common.SystemStatsResponse
	err := n.createRequest(n.client, "GET", "stats/system", &common.Empty{}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetBackendStats() (*common.BackendStatsResponse, error) {
	var stats common.BackendStatsResponse
	err := n.createRequest(n.client, "GET", "stats/backend", &common.Empty{}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetStats(reset bool, name string, statType common.StatType) (*common.StatResponse, error) {
	var stats common.StatResponse
	if err := n.createRequest(n.client, "GET", "stats", &common.StatRequest{Reset_: reset, Name: name, Type: statType}, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetUserOnlineStat(email string) (*common.OnlineStatResponse, error) {
	var stats common.OnlineStatResponse
	err := n.createRequest(n.client, "GET", "stats/user/online", &common.StatRequest{Name: email}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (n *Node) GetUserOnlineIpList(email string) (*common.StatsOnlineIpListResponse, error) {
	var stats common.StatsOnlineIpListResponse
	err := n.createRequest(n.client, "GET", "stats/user/online_ip", &common.StatRequest{Name: email}, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
