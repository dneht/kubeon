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

package cluster

import (
	"github.com/dneht/kubeon/pkg/onutil/log"
	"testing"
)

func TestMain(m *testing.M) {
	log.Init(6)
	m.Run()
}

func TestNodeListSort_Run(t *testing.T) {
	list := NodeList{
		&Node{
			IPv4:  "1.1.1.1",
			Order: 2,
		},
		&Node{
			IPv4:  "9.1.1.1",
			Order: 9,
		},
		&Node{
			IPv4:  "0.1.1.1",
			Order: 0,
		},
		&Node{
			IPv4:  "3.1.1.1",
			Order: 3,
		},
	}
	for _, n := range SortNodeList(list) {
		t.Log(n.IPv4)
	}
}
