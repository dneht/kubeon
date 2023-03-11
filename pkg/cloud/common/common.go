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
	"github.com/dneht/kubeon/pkg/cluster"
	"os"
)

func GetAccessKeyFromEnv(provider string) (string, string) {
	id, iok := os.LookupEnv("CLOUD_SECRET_ID")
	if !iok {
		return "", ""
	}
	key, sok := os.LookupEnv("CLOUD_SECRET_KEY")
	if !sok {
		return "", ""
	}
	return id, key
}

func GetRouterNamePrefix() string {
	name := "unknown"
	if nil != cluster.Current() {
		name = cluster.Current().Name
	}
	return "kubeon-" + name + "-pod-host-"
}

func GetRouterDescPrefix() string {
	name := "unknown"
	if nil != cluster.Current() {
		name = cluster.Current().Name
	}
	return "Routing rules set by kubeon for vpc, cluster: " + name + ", node: "
}

func CopyNodeInfoList(nodeInfos map[string]*cluster.NodeCloudInfo) map[string]*cluster.NodeCloudInfo {
	nodeNews := make(map[string]*cluster.NodeCloudInfo, len(nodeInfos))
	for key, value := range nodeInfos {
		nodeNews[key] = &cluster.NodeCloudInfo{
			Name:       value.Name,
			Desc:       value.Desc,
			IP:         value.IP,
			CIDR:       value.CIDR,
			InstanceId: value.InstanceId,
			EntryId:    value.EntryId,
			RouterId:   value.RouterId,
		}
	}
	return nodeNews
}
