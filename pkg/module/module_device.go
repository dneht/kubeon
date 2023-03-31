/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package module

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"k8s.io/klog/v2"
)

func LabelDevice() {
	LabelNvidia()
}

func LabelNvidia() {
	current := cluster.Current()
	for _, node := range cluster.CurrentNodes() {
		if current.UseNvidia && node.HasNvidia {
			err := action.KubectlLabelRole(node.Hostname, "nvidia.com/gpu.present=yes")
			if nil != err {
				klog.Warningf("Label[nvidia.com/gpu.present=yes] on %s failed: %v, please set it manually using [kubectl label nodes %s nvidia.com/gpu.present=yes]", node.Hostname, err, node.Hostname)
			}
		}
	}
}
