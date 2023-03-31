/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package cluster

type NodeCloudInfo struct {
	Name       string
	Desc       string
	IP         string
	CIDR       string
	InstanceId string
	EntryId    string
	RouterId   uint64
}
