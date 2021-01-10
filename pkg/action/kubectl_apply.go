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
	"bytes"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/execute"
)

func KubectlApplyPath(path string) error {
	return execute.NewLocalCmd("kubectl", "apply",
		"--kubeconfig="+cluster.Current().AdminConfigPath, "-f", path).RunWithEcho()
}

func KubectlApplyData(content []byte) error {
	cmd := execute.NewLocalCmd("kubectl", "apply", "--kubeconfig="+cluster.Current().AdminConfigPath, "-f", "-")
	cmd.Stdin(bytes.NewBuffer(content))
	return cmd.RunWithEcho()
}

func KubectlDeleteData(content []byte) error {
	cmd := execute.NewLocalCmd("kubectl", "delete", "--kubeconfig="+cluster.Current().AdminConfigPath, "-f", "-")
	cmd.Stdin(bytes.NewBuffer(content))
	return cmd.RunWithEcho()
}
