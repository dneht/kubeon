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

package define

import "sort"

var (
	supportVersions = map[string]uint{
		K8S_1_19_4.Full: K8S_1_19_4.Number,
		K8S_1_19_5.Full: K8S_1_19_5.Number,
		K8S_1_19_6.Full: K8S_1_19_6.Number,
		K8S_1_19_7.Full: K8S_1_19_7.Number,
		K8S_1_19_8.Full: K8S_1_19_8.Number,
		K8S_1_19_9.Full: K8S_1_19_9.Number,
		K8S_1_20_1.Full: K8S_1_20_1.Number,
		K8S_1_20_2.Full: K8S_1_20_2.Number,
		K8S_1_20_3.Full: K8S_1_20_3.Number,
		K8S_1_20_4.Full: K8S_1_20_4.Number,
		K8S_1_20_5.Full: K8S_1_20_5.Number,
	}
	supportProxyModes = map[string]bool{
		IPVSProxy:     true,
		IPTablesProxy: true,
	}
	supportRuntimes = map[string]bool{
		DockerRuntime:     true,
		ContainerdRuntime: true,
	}
	supportNetworks = map[string]bool{
		NoneNetwork:   true,
		CalicoNetwork: true,
	}
	supportCalicoModes = map[string]bool{
		CalicoIPIP:  true,
		CalicoVXLan: true,
	}
)

func SupportVersionList() []string {
	versions := make(map[uint]string, len(supportVersions))
	for k, v := range supportVersions {
		versions[v] = k
	}
	keys := make([]int, 0, len(versions))
	for key := range versions {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)
	supports := make([]string, 0, len(keys))
	for _, key := range keys {
		supports = append(supports, versions[uint(key)])
	}
	return supports
}

func IsSupportVersion(version string) bool {
	_, ok := supportVersions[version]
	return ok
}

func SupportRuntimeList() []string {
	return getMapKeys(supportRuntimes)
}

func IsSupportRuntime(cri string) bool {
	_, ok := supportRuntimes[cri]
	return ok
}

func SupportNetworkList() []string {
	return getMapKeys(supportNetworks)
}

func IsSupportNetwork(cni string) bool {
	_, ok := supportNetworks[cni]
	return ok
}

func SupportProxyModes() []string {
	return getMapKeys(supportProxyModes)
}

func IsSupportProxyMode(mode string) bool {
	_, ok := supportProxyModes[mode]
	return ok
}

func SupportCalicoModes() []string {
	return getMapKeys(supportCalicoModes)
}

func IsSupportCalicoMode(mode string) bool {
	_, ok := supportCalicoModes[mode]
	return ok
}

func getMapKeys(input map[string]bool) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	return keys
}
