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
