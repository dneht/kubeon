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

var (
	supportVersions = []*RngVersion{
		K8S_1_19_x,
		K8S_1_20_x,
		K8S_1_21_x,
		K8S_1_22_x,
		K8S_1_23_x,
		K8S_1_24_x,
		K8S_1_25_x,
		K8S_1_26_x,
	}
	supportProxyModes = map[string]bool{
		IPVSProxy:     true,
		IPTablesProxy: true,
		BPFProxy:      true,
		CalicoNetwork: true,
		CiliumNetwork: true,
	}
	supportRuntimes = map[string]bool{
		DockerRuntime:     true,
		ContainerdRuntime: true,
	}
	supportNetworks = map[string]bool{
		NoneNetwork:   true,
		BPFProxy:      true,
		CalicoNetwork: true,
		CiliumNetwork: true,
	}
	supportIngresses = map[string]bool{
		NoneIngress:    true,
		ContourIngress: true,
		IstioIngress:   true,
	}
	supportCiliumPolicyMode = map[string]bool{
		CiliumPolicyDefault: true,
		CiliumPolicyAlways:  true,
		CiliumPolicyNever:   true,
	}
)

func SupportVersionList() []string {
	supports := make([]string, 0, len(supportVersions))
	for _, ver := range supportVersions {
		supports = append(supports, ver.String())
	}
	return supports
}

func IsSupportVersion(version string) bool {
	for _, ver := range supportVersions {
		if ver.Contain(version) {
			return true
		}
	}
	return false
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

func IsSupportIngress(ic string) bool {
	_, ok := supportIngresses[ic]
	return ok
}

func SupportProxyModes() []string {
	return getMapKeys(supportProxyModes)
}

func IsSupportProxyMode(mode string) bool {
	_, ok := supportProxyModes[mode]
	return ok
}

func IsSupportCiliumPolicyMode(mode string) bool {
	_, ok := supportCiliumPolicyMode[mode]
	return ok
}

func getMapKeys(input map[string]bool) []string {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	return keys
}
