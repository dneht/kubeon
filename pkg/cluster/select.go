/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package cluster

import (
	"github.com/pkg/errors"
	"strings"
)

const (
	AllNode           = "@all"
	AllControlNode    = "@cp*"
	NowBootNode       = "@cp1"
	NoBootControlNode = "@cpn"
	AllJoinNode       = "@j*"
	AllWorkerNode     = "@w*"
)

func SelectNodes(nodeSelector string) (nodes NodeList, err error) {
	if nil == current {
		return nil, errors.New("cluster is not init, please check context")
	}
	if strings.HasPrefix(nodeSelector, "@") {
		switch strings.ToLower(nodeSelector) {
		case AllNode:
			return current.AllNodes, nil
		case AllControlNode:
			return current.ControlPlanes, nil
		case NowBootNode:
			if len(current.ControlPlanes) == 0 {
				return nil, nil
			}
			return toNodeList(BootstrapNode()), nil
		case NoBootControlNode:
			if len(current.ControlPlanes) <= 1 {
				return nil, nil
			}
			return current.ControlPlanes[1:], nil
		case AllJoinNode:
			if len(current.AllNodes) <= 1 {
				return nil, nil
			}
			return current.AllNodes[1:], nil
		case AllWorkerNode:
			return current.Workers, nil
		default:
			return nil, errors.Errorf("Invalid node selector %q. Use one of [@all, @cp*, @cp1, @cpn, @w*, @lb, @etcd]", nodeSelector)
		}
	} else {
		idx := strings.Index(nodeSelector, "=")
		if idx > 0 {
			selectType := nodeSelector[0:idx]
			selectVals := nodeSelector[idx+1:]
			selectValArr := strings.Split(selectVals, ",")
			nodeList := NodeList{}
			switch strings.ToLower(selectType) {
			case "ip":
				for _, n := range current.AllNodes {
					for _, s := range selectValArr {
						if strings.EqualFold(s, n.IPv4) {
							nodeList = append(nodeList, n)
						}
					}
				}
				return nodeList, nil
			case "n", "name":
				for _, n := range current.AllNodes {
					for _, s := range selectValArr {
						if strings.EqualFold(s, n.Hostname) {
							nodeList = append(nodeList, n)
						}
					}
				}
				return nodeList, nil
			default:
				return nil, errors.Errorf("Invalid node selector %q. Use one of [ip=ipv4, name=hostname]", nodeSelector)
			}
		}
	}
	for _, n := range current.AllNodes {
		if strings.EqualFold(nodeSelector, n.Addr()) {
			return toNodeList(n), nil
		}
	}
	return nil, nil
}

func BootstrapNode() *Node {
	return current.AllNodes[0]
}

func JoinsNode() NodeList {
	return current.AllNodes[1:]
}

func SelectNodeByIP(ipv4 string) *Node {
	for _, n := range current.AllNodes {
		if strings.EqualFold(ipv4, n.IPv4) {
			return n
		}
	}
	return nil
}

func SelectNodeByName(name string) *Node {
	for _, n := range current.AllNodes {
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
