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

const patchCorednsJson = "[{\"op\":\"add\",\"path\":\"/rules/-\",\"value\":{\"apiGroups\":[\"discovery.k8s.io\"],\"resources\":[\"endpointslices\"],\"verbs\":[\"list\",\"watch\"]}}]"

func KubectlPatchCorednsRole() error {
	output, err := KubectlGetResult(
		"clusterrole",
		"system:coredns",
		"-o=jsonpath='{.rules[-1].apiGroups}'",
	)
	if nil != err {
		return err
	}
	if output == "'[\"\"]'" {
		return execute.NewLocalCmd("kubectl",
			"patch", "clusterrole", "system:coredns",
			"--type=json", "--patch="+patchCorednsJson,
			"--kubeconfig="+cluster.Current().AdminConfigPath,
		).RunWithEcho()
	}
	return nil
}
