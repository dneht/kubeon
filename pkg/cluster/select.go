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

package cluster

import (
	"github.com/pkg/errors"
	"strings"
)

const (
	AllNode        = "@all"
	AllControlNode = "@cp*"
	NowBootNode    = "@cp1"
	AllJoinNode    = "@j*"
	AllWorkerNode  = "@w*"
)

func SelectNodes(nodeSelector string) (nodes NodeList, err error) {
	if nil == current {
		return nil, errors.New("cluster is not init, please check context")
	}
	if strings.HasPrefix(nodeSelector, "@") {
		switch strings.ToLower(nodeSelector) {
		case "@all":
			return currentNodes, nil
		case "@cp*":
			return current.ControlPlanes, nil
		case "@cp1":
			if len(current.ControlPlanes) == 0 {
				return nil, nil
			}
			return toNodeList(BootstrapNode()), nil
		case "@cpn":
			if len(current.ControlPlanes) <= 1 {
				return nil, nil
			}
			return current.ControlPlanes[1:], nil
		case "@j*":
			if len(currentNodes) <= 1 {
				return nil, nil
			}
			return currentNodes[1:], nil
		case "@w*":
			return current.Workers, nil
		default:
			return nil, errors.Errorf("Invalid node selector %q. Use one of [@all, @cp*, @cp1, @cpn, @w*, @lb, @etcd]", nodeSelector)
		}
	} else {
		idx := strings.Index(nodeSelector, "=")
		if idx > 0 {
			selectType := nodeSelector[0:idx]
			selectVal := nodeSelector[idx+1:]
			switch strings.ToLower(selectType) {
			case "ip":
				for _, n := range currentNodes {
					if strings.EqualFold(selectVal, n.IPv4) {
						return toNodeList(n), nil
					}
				}
			case "name":
				for _, n := range currentNodes {
					if strings.EqualFold(selectVal, n.Hostname) {
						return toNodeList(n), nil
					}
				}
			default:
				return nil, errors.Errorf("Invalid node selector %q. Use one of [ip=ipv4, name=hostname]", nodeSelector)
			}
		}
	}
	for _, n := range currentNodes {
		if strings.EqualFold(nodeSelector, n.Addr()) {
			return toNodeList(n), nil
		}
	}
	return nil, nil
}

func BootstrapNode() *Node {
	return currentNodes[0]
}

func JoinsNode() NodeList {
	return currentNodes[1:]
}

func SelectNodeByIP(ipv4 string) *Node {
	for _, n := range currentNodes {
		if strings.EqualFold(ipv4, n.IPv4) {
			return n
		}
	}
	return nil
}

func SelectNodeByName(name string) *Node {
	for _, n := range currentNodes {
		if strings.EqualFold(name, n.Home) {
			return n
		}
	}
	return nil
}

func GetMasterFromList(nodes NodeList) NodeList {
	list := make(NodeList, 0, len(nodes)/2+1)
	for _, node := range nodes {
		if node.IsControlPlane() {
			list = append(list, node)
		}
	}
	return list
}

func toNodeList(node *Node) NodeList {
	if node != nil {
		return NodeList{node}
	}
	return nil
}
