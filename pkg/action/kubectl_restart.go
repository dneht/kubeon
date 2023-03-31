/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/execute"
	"k8s.io/klog/v2"
	"strings"
)

func RestartCalicoNodeForce() {
	RestartDaemonSetForce("calico-node")
}

func RestartKubeProxyForce() {
	RestartDaemonSetForce("kube-proxy")
}

func RestartDaemonSetForce(label string) {
	lines, err := execute.NewLocalCmd("kubectl",
		"get", "pods",
		"--namespace=kube-system",
		"--selector=k8s-app="+label,
		"--kubeconfig="+cluster.Current().AdminConfigPath).RunAndCapture()
	if nil != err {
		klog.Warningf("Get %s pods failed: %v", label, err)
		return
	}
	for _, one := range lines {
		one = strings.TrimSpace(one)
		if strings.HasPrefix(one, label) {
			err = execute.NewLocalCmd("kubectl",
				"delete", "pods", strings.TrimSpace(strings.Split(one, " ")[0]),
				"--namespace=kube-system",
				"--kubeconfig="+cluster.Current().AdminConfigPath).RunWithEcho()
			if nil != err {
				klog.Warningf("Delete %s pod[%s] failed: %v", label, one, err)
			}
		}
	}
}
