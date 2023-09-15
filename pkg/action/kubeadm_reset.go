/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"k8s.io/klog/v2"
)

func KubeadmResetOne(node *cluster.Node, delete, force bool) {
	var err error
	current := cluster.Current()
	if force {
		err = node.RunCmd("systemctl", "stop", current.RuntimeMode, "--force")
		if nil != err {
			klog.Warningf("%s restart failed: %v", current.RuntimeMode, err)
		}
	}
	err = node.RunCmd("kubeadm", "reset", "--force", fmt.Sprintf("--v=%d", log.Level()))
	if nil != err {
		klog.Warningf("Kubeadm reset failed: %v", err)
	}
	err = node.Rm("/etc/cni/net.d")
	if nil != err {
		klog.Warningf("Remove cni config failed: %v", err)
	}
	_ = node.Rm("/etc/kubernetes")
	_ = node.Rm("/etc/kubeadm.yaml")
	if current.ProxyMode == define.IPVSProxy {
		err = node.RunCmd("ipvsadm", "--clear")
		if nil != err {
			klog.Warningf("Clean ipvs rules failed: %v", err)
		}
	} else if current.ProxyMode == define.IPTablesProxy {
		klog.Warningf("Please clean the iptables rules yourself")
	}
	if delete {
		err = KubectlDeleteNodeForce(node.Hostname)
		if nil != err {
			klog.Warningf("Delete node[%s] failed: %v", node.Addr(), err)
		}
	}
}

func KubeadmResetList(list cluster.NodeList, delete, force bool) {
	for _, node := range list {
		if node.IsWorker() {
			KubeadmResetOne(node, delete, false)
		}
	}
	var boot *cluster.Node = nil
	for _, node := range list {
		if node.IsBootstrap() {
			boot = node
			continue
		}
		if node.IsControlPlane() {
			KubeadmResetOne(node, delete, false)
		}
	}
	if nil != boot {
		KubeadmResetOne(boot, delete, force)
	}
}
