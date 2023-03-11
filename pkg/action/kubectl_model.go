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

type NodeInfo struct {
	Spec   *NodeSpecInfo   `json:"spec"`
	Status *NodeStatusInfo `json:"status"`
}

type NodeSpecInfo struct {
	PodCIDR  string   `json:"podCIDR"`
	PodCIDRs []string `json:"podCIDRs"`
}

type NodeStatusInfo struct {
	NodeInfo  *NodeBaseInfo  `json:"nodeInfo"`
	Addresses []NodeAddrInfo `json:"addresses"`
}

type NodeBaseInfo struct {
	Architecture            string `json:"architecture"`
	BootID                  string `json:"bootID"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KernelVersion           string `json:"kernelVersion"`
	KubeProxyVersion        string `json:"kubeProxyVersion"`
	KubeletVersion          string `json:"kubeletVersion"`
	MachineID               string `json:"machineID"`
	OperatingSystem         string `json:"operatingSystem"`
	OsImage                 string `json:"osImage"`
	SystemUUID              string `json:"systemUUID"`
}

type NodeAddrInfo struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}
