/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package module

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/release"
	"k8s.io/klog/v2"
)

func AllUninstall(nodes cluster.NodeList, isDestroy bool) (err error) {
	return removePackage(nodes, isDestroy)
}

func removePackage(nodes cluster.NodeList, isDestroy bool) (err error) {
	for _, node := range nodes {
		klog.V(1).Infof("Start uninstall [%s] on [%s]", define.KubeletModule, node.Addr())
		err = uninstallOnNode(node, define.KubeletModule)
		if nil != err {
			klog.Warningf("Uninstall module %s error: %s", define.KubeletModule, err)
		} else {
			klog.V(4).Infof("Uninstall module %s success", define.KubeletModule)
		}

		klog.V(1).Infof("Start uninstall [%s] on [%s]", cluster.Current().RuntimeMode, node.Addr())
		if cluster.Current().RuntimeMode == define.ContainerdRuntime {
			err = uninstallOnNode(node, define.ContainerdRuntime)
			if nil != err {
				klog.Warningf("Uninstall module %s error: %s", define.ContainerdRuntime, err)
			} else {
				klog.V(4).Infof("Uninstall module %s success", define.ContainerdRuntime)
			}
		} else {
			err = uninstallOnNode(node, define.DockerRuntime)
			if nil != err {
				klog.Warningf("Uninstall module %s error: %s", define.DockerRuntime, err)
			} else {
				klog.V(4).Infof("Uninstall module %s success", define.DockerRuntime)
			}
		}
		err = uninstallOnNode(node, define.NetworkPlugin)
		if nil != err {
			klog.Warningf("Uninstall module %s error: %s", define.NetworkPlugin, err)
		} else {
			klog.V(4).Infof("Uninstall module %s success", define.NetworkPlugin)
		}
		uninstallScript(node)
		if onutil.IsLocalIPv4(node.IPv4) {
			release.ReinstallLocal(cluster.CurrentResource())
		} else {
			_ = node.Rm(node.GetResource().BaseDir)
			_ = node.Rm(node.Home + "/.kube")
			_ = node.Rm("/usr/local/bin/kubeon")
		}
	}
	return nil
}

func uninstallScript(node *cluster.Node) {
	current := cluster.Current()
	installMode := "online"
	if current.IsOffline {
		installMode = "offline"
	}
	proxyMode := current.ProxyMode
	klog.V(1).Infof("Start final uninstall on [%s], %s, proxy=%s", node.Addr(), installMode, proxyMode)
	err := node.RunCmd("bash", node.GetResource().ScriptDir+"/prepare.sh",
		"delete", installMode, proxyMode)
	if nil != err {
		klog.Warningf("Final uninstall on [%s] failed", node.Addr())
	}
}
