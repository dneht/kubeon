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

package module

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/dneht/kubeon/pkg/release"
)

func AllUninstall(nodes cluster.NodeList, isDestroy bool) (err error) {
	return removePackage(nodes, isDestroy)
}

func removePackage(nodes cluster.NodeList, isDestroy bool) (err error) {
	for _, node := range nodes {
		log.Infof("start uninstall [%s] on [%s]", define.KubeletModule, node.Addr())
		err = uninstallOnNode(node, define.KubeletModule)
		if nil != err {
			log.Warnf("uninstall module %s error: %s", define.KubeletModule, err)
		} else {
			log.Debugf("uninstall module %s success", define.KubeletModule)
		}

		log.Infof("start uninstall [%s] on [%s]", cluster.Current().RuntimeMode, node.Addr())
		if cluster.Current().RuntimeMode == define.ContainerdRuntime {
			err = uninstallOnNode(node, define.ContainerdRuntime)
			if nil != err {
				log.Warnf("uninstall module %s error: %s", define.ContainerdRuntime, err)
			} else {
				log.Debugf("uninstall module %s success", define.ContainerdRuntime)
			}
		} else {
			err = uninstallOnNode(node, define.DockerRuntime)
			if nil != err {
				log.Warnf("uninstall module %s error: %s", define.DockerRuntime, err)
			} else {
				log.Debugf("uninstall module %s success", define.DockerRuntime)
			}
		}
		err = uninstallOnNode(node, define.NetworkPlugin)
		if nil != err {
			log.Warnf("uninstall module %s error: %s", define.NetworkPlugin, err)
		} else {
			log.Debugf("uninstall module %s success", define.NetworkPlugin)
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
	log.Infof("start final uninstall on [%s], %s, proxy=%s", node.Addr(), installMode, proxyMode)
	err := node.RunCmd("bash", node.GetResource().ScriptDir+"/prepare.sh",
		"delete", installMode, proxyMode)
	if nil != err {
		log.Warnf("final uninstall on [%s] failed", node.Addr())
	}
}
