/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package common

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
)

func AllNodeCloudInfo() (map[string]*cluster.NodeCloudInfo, error) {
	allNodes := cluster.CurrentNodes()
	nodeInfoList := make(map[string]*cluster.NodeCloudInfo, len(allNodes))
	for _, node := range allNodes {
		podCIDR, err := action.GetNodePodCIDR(node)
		if nil != err {
			return nil, err
		}
		nodeInfoList[node.IPv4] = &cluster.NodeCloudInfo{
			Name: node.Hostname,
			IP:   node.IPv4,
			CIDR: podCIDR,
		}
	}
	return nodeInfoList, nil
}
