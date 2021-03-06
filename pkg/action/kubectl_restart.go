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

package action

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil/log"
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
		log.Warnf("get %s pods error: %s", label, err)
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
				log.Warnf("delete %s pod[%s] error: %s", label, one, err)
			}
		}
	}
}
