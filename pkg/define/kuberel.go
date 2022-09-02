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
	K8S_1_15_0, _ = NewStdVersion("v1.15.0")
	K8S_1_19_0, _ = NewStdVersion("v1.19.0")
	K8S_1_20_0, _ = NewStdVersion("v1.20.0")
	K8S_1_21_0, _ = NewStdVersion("v1.21.0")
	K8S_1_22_0, _ = NewStdVersion("v1.22.0")
	K8S_1_23_0, _ = NewStdVersion("v1.23.0")
	K8S_1_19_x, _ = NewRngVersion("v1.19.4", "v1.19.16")
	K8S_1_20_x, _ = NewRngVersion("v1.20.1", "v1.20.15")
	K8S_1_21_x, _ = NewRngVersion("v1.21.1", "v1.21.14")
	K8S_1_22_x, _ = NewRngVersion("v1.22.1", "v1.22.13")
	K8S_1_23_x, _ = NewRngVersion("v1.23.1", "v1.23.10")
	K8S_1_24_x, _ = NewRngVersion("v1.24.1", "v1.24.4")
	ETCD_3_4_0, _ = NewStdVersion("3.4.0")
)

const (
	KubeletConfigApiB1 = "v1beta1"
	KubeadmConfigApiB2 = "v1beta2"
	KubeadmConfigApiB3 = "v1beta3"
)

const (
	ImagesPackage = "images"
	PausePackage  = "pause"
)

const (
	KubeletModule = "kubelet"
	KubeadmModule = "kubeadm"
	BinaryModule  = "binary"
	OnlineModule  = "online"
	OfflineModule = "offline"
)

const InstallScript = "script"

const (
	HealthzReader    = "healthz"
	LocalHaproxy     = "haproxy"
	ApiserverUpdater = "updater"
	ApiserverStartup = "startup"
	ApiserverService = "apiserver-startup"
)

func (v *StdVersion) IsSupportPatch() bool {
	return v.GreaterEqual(K8S_1_19_0)
}

func (v *StdVersion) IsSupportContour() bool {
	return v.GreaterThen(K8S_1_21_0)
}

func (v *StdVersion) IsSupportNvidia() bool {
	return v.GreaterThen(K8S_1_22_0)
}

func (v *StdVersion) IsSupportKata() bool {
	return v.GreaterThen(K8S_1_22_0)
}
