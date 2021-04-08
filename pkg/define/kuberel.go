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
	K8S_1_19_4, _ = NewStdVersion("v1.19.4")
	K8S_1_19_5, _ = NewStdVersion("v1.19.5")
	K8S_1_19_6, _ = NewStdVersion("v1.19.6")
	K8S_1_19_7, _ = NewStdVersion("v1.19.7")
	K8S_1_19_8, _ = NewStdVersion("v1.19.8")
	K8S_1_19_9, _ = NewStdVersion("v1.19.9")
	K8S_1_20_1, _ = NewStdVersion("v1.20.1")
	K8S_1_20_2, _ = NewStdVersion("v1.20.2")
	K8S_1_20_3, _ = NewStdVersion("v1.20.3")
	K8S_1_20_4, _ = NewStdVersion("v1.20.4")
	K8S_1_20_5, _ = NewStdVersion("v1.20.5")
	ETCD_3_4_0, _ = NewStdVersion("3.4.0")
)

const ImagesPackage = "images"

const (
	KubeletModule = "kubelet"

	BinaryModule = "binary"

	OfflineModule = "offline"
)

const (
	HelmTool = "helm"
)

const (
	HealthzReader = "healthz"

	LocalHaproxy = "haproxy"

	ApiserverUpdater = "updater"

	ApiserverStartup = "startup"

	ApiserverService = "apiserver-startup"
)
