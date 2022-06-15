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
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/execute"
)

func KubectlGetQuiet(args ...string) (string, error) {
	return kubectlGet(args).Quiet().RunAndResult()
}

func KubectlGetResult(args ...string) (string, error) {
	return kubectlGet(args).RunAndResult()
}

func KubectlGetOutput(args ...string) error {
	return kubectlGet(args).RunWithEcho()
}

func kubectlGet(args []string) *execute.LocalCmd {
	getArgs := []string{
		"get",
		"--kubeconfig=" + cluster.Current().AdminConfigPath,
	}
	getArgs = append(getArgs, args...)
	return execute.NewLocalCmd("kubectl", getArgs...)
}

func KubectlAuthGetResult(args ...string) (string, error) {
	getArgs := []string{
		"auth",
		"can-i",
		"get",
		"--kubeconfig=" + cluster.Current().AdminConfigPath,
	}
	getArgs = append(getArgs, args...)
	return execute.NewLocalCmd("kubectl", getArgs...).RunAndResult()
}
