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
)

func KubectlRemoveMasterTaint(name string) {
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes", name,
		"node-role.kubernetes.io/master-",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes", name,
		"node-role.kubernetes.io/control-plane-",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
}

func KubectlRemoveAllMasterTaint() {
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes",
		"node-role.kubernetes.io/master-",
		"--all",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes",
		"node-role.kubernetes.io/control-plane-",
		"--all",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
}
