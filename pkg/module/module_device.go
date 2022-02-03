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
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"k8s.io/klog/v2"
)

func LabelDevice() {
	LabelNvidia()
}

func LabelNvidia() {
	current := cluster.Current()
	for _, node := range cluster.CurrentNodes() {
		if current.UseNvidia && node.HasNvidia {
			err := action.KubectlLabelRole(node.Hostname, "nvidia.com/gpu.present=yes")
			if nil != err {
				klog.Warningf("Label[nvidia.com/gpu.present=yes] on %s failed: %v, please set it manually using [kubectl label nodes %s nvidia.com/gpu.present=yes]", node.Hostname, err, node.Hostname)
			}
		}
	}
}
