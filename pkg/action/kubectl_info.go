/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"encoding/json"
	"github.com/dneht/kubeon/pkg/cluster"
	"strings"
)

func GetNodePodCIDR(node *cluster.Node) (string, error) {
	result, err := KubectlGetResult(
		"nodes", node.Hostname, "-o=jsonpath='{.spec.podCIDR}'",
	)
	if nil != err {
		return "", err
	}
	return strings.ReplaceAll(result, "'", ""), nil
}

func GetNodeInfo(node *cluster.Node) (*NodeInfo, error) {
	output, err := KubectlGetQuiet(
		"nodes", node.Hostname, "-o=json",
	)
	if nil != err {
		return nil, err
	}
	var info NodeInfo
	err = json.Unmarshal([]byte(output), &info)
	if nil != err {
		return nil, err
	}
	return &info, nil
}
