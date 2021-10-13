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
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
)

func deleteLocalFlag(version *define.StdVersion) string {
	if version.GreaterEqual(define.K8S_1_20_0) {
		return "--delete-emptydir-data"
	} else {
		return "--delete-local-data"
	}
}

func KubectlDrainNode(name string, version *define.StdVersion) (err error) {
	return execute.NewLocalCmd("kubectl",
		"drain", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath,
		"--ignore-daemonsets", deleteLocalFlag(version)).RunWithEcho()
}

func KubectlDrainNodeForce(name string, version *define.StdVersion) (err error) {
	return execute.NewLocalCmd("kubectl",
		"drain", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath,
		"--force", "--ignore-daemonsets", deleteLocalFlag(version)).RunWithEcho()
}

func KubectlUncordonNode(name string) (err error) {
	return execute.NewLocalCmd("kubectl",
		"uncordon", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath).RunWithEcho()
}
