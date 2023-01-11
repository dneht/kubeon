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

package define

type ComponentVersion struct {
	Kubernetes string `json:"kubernetes"`
	Pause      string `json:"pause"`
	Etcd       string `json:"etcd"`
	Coredns    string `json:"coredns"`
	Crictl     string `json:"crictl,omitempty"`
	Runc       string `json:"runc,omitempty"`
	Containerd string `json:"containerd,omitempty"`
	Docker     string `json:"docker,omitempty"`
	Nvidia     string `json:"nvidia,omitempty"`
	Kata       string `json:"kata,omitempty"`
	Network    string `json:"cni,omitempty"`
	Calico     string `json:"calico,omitempty"`
	Cilium     string `json:"cilium,omitempty"`
	Contour    string `json:"contour,omitempty"`
	Istio      string `json:"istio,omitempty"`
	Haproxy    string `json:"haproxy,omitempty"`
	Kruise     string `json:"kruise,omitempty"`
	Offline    string `json:"offline"`
}

var SupportComponentFull = map[string]*ComponentVersion{}
