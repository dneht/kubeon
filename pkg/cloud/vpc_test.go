/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package cloud

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"testing"
)

func TestModifyRouter(t *testing.T) {
	info := map[string]*cluster.NodeCloudInfo{
		"172.16.0.10": {
			Name: "test30",
			IP:   "172.16.0.10",
			CIDR: "10.96.30.0/24",
		},
		"172.16.0.11": {
			Name: "test31",
			IP:   "172.16.0.11",
			CIDR: "10.96.31.0/24",
		},
		"172.16.0.12": {
			Name: "test32",
			IP:   "172.16.0.12",
			CIDR: "10.96.32.0/24",
		},
	}
	err := ModifyRouter("xxx", &cluster.CloudConf{
		Endpoint:       "x-shanghai",
		RouterTableIds: []string{"xxx-xxxx"},
	}, info)
	t.Logf("%v\n", info)
	if err != nil {
		t.Error(err)
	}
}
