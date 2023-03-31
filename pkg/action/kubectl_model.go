/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
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
