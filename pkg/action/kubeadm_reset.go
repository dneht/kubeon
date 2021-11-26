/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Full 2.0 (the "License");
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

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil/log"
)

func KubeadmResetOne(node *cluster.Node) {
	err := node.RunCmd("kubeadm", "reset", "--force", fmt.Sprintf("--v=%d", log.Level()))
	if nil != err {
		log.Warnf("kubeadm reset failed: %v", err)
	}
	err = node.Rm("/etc/cni/net.d")
	if nil != err {
		log.Warnf("remove cni config failed: %v", err)
	}
	if cluster.Current().ProxyMode == define.IPVSProxy {
		err = node.RunCmd("ipvsadm", "--clear")
		if nil != err {
			log.Warnf("clean ipvs rules failed: %v", err)
		}
	} else if cluster.Current().ProxyMode == define.IPTablesProxy {
		log.Warnf("please clean the iptables rules yourself")
	}
}

func KubeadmResetList(list cluster.NodeList) {
	for _, node := range list {
		if node.IsWorker() {
			KubeadmResetOne(node)
		}
	}
	var boot *cluster.Node = nil
	for _, node := range list {
		if node.IsBootstrap() {
			boot = node
			continue
		}
		if node.IsControlPlane() {
			KubeadmResetOne(node)
		}
	}
	if nil != boot {
		KubeadmResetOne(boot)
	}
}
