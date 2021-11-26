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
	"github.com/dneht/kubeon/pkg/onutil/log"
)

func IPVSAdmClear(nodes cluster.NodeList) {
	var err error
	for _, node := range nodes {
		err = node.RunCmd("ipvsadm", "-C")
		if nil != err {
			log.Warnf("clear ipvs rule failed: %v", err)
		}
	}
}
