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

package cloud

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"testing"
)

func TestModifyRouter(t *testing.T) {
	info := map[string]*cluster.NodeCloudInfo{
		"172.16.0.10": {
			Name: "test30",
			IP:   "172.16.0.10",
			CIDR: "10.96.30.0/24",
		},
		"172.16.0.11": {
			Name: "test31",
			IP:   "172.16.0.11",
			CIDR: "10.96.31.0/24",
		},
		"172.16.0.12": {
			Name: "test32",
			IP:   "172.16.0.12",
			CIDR: "10.96.32.0/24",
		},
	}
	err := ModifyRouter("xxx", &cluster.CloudConf{
		Endpoint:       "x-shanghai",
		RouterTableIds: []string{"xxx-xxxx"},
	}, info)
	t.Logf("%v\n", info)
	if err != nil {
		t.Error(err)
	}
}
