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

package create

import (
	"github.com/dneht/kubeon/pkg/define"
	"k8s.io/klog/v2"
)

func checkSupport(flags *flagpole, clusterVersion string) bool {
	var isSupport bool
	isSupport = define.IsSupportVersion(clusterVersion)
	if !isSupport {
		klog.Errorf("[check] Input version[%s] not support", clusterVersion)
		return false
	}
	isSupport = define.IsSupportRuntime(flags.InputCRIMode)
	if !isSupport {
		klog.Errorf("[check] Input cri[%s] not support", flags.InputCRIMode)
		return false
	}
	isSupport = define.IsSupportNetwork(flags.InputCNIMode)
	if !isSupport {
		klog.Errorf("[check] Input cni[%s] not support", flags.InputCNIMode)
		return false
	}
	isSupport = define.IsSupportIngress(flags.InputICMode)
	if !isSupport {
		klog.Errorf("[check] Input ic[%s] not support", flags.InputICMode)
		return false
	}
	isSupport = define.IsSupportProxyMode(flags.InputProxyMode)
	if !isSupport {
		klog.Errorf("[check] Input proxy mode[%s] not support", flags.InputProxyMode)
		return false
	}
	isSupport = define.IsSupportCalicoMode(flags.CalicoMode)
	if !isSupport {
		klog.Errorf("[check] Input calico mode[%s] not support", flags.CalicoMode)
		return false
	}
	return true
}
