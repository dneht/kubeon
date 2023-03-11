/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
