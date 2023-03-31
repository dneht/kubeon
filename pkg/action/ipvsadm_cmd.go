/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"k8s.io/klog/v2"
)

func IPVSAdmClear(nodes cluster.NodeList) {
	var err error
	for _, node := range nodes {
		err = node.RunCmd("ipvsadm", "-C")
		if nil != err {
			klog.Warningf("Clear ipvs rule failed: %v", err)
		}
	}
}
