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
		err = KubectlDeleteNode(node.Hostname)
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
