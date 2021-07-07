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
	K8S_1_19_x, _ = NewRngVersion("v1.19.4", "v1.19.12")
	K8S_1_20_x, _ = NewRngVersion("v1.20.1", "v1.20.8")
	K8S_1_21_x, _ = NewRngVersion("v1.21.1", "v1.21.2")
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
