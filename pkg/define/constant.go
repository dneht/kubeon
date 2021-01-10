/*
Copyright 2019 The Kubernetes Authors.

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

import "github.com/dneht/kubeon/pkg/onutil"

const (
	AppVersion = "0.7.5"

	MirrorRegistry = "registry.cn-hangzhou.aliyuncs.com"
)

var (
	DefaultKubeVersion = K8S_1_19_4

	AppBaseDir = onutil.BaseDir()

	AppConfDir = AppBaseDir + "/conf"

	AppTplDir = AppBaseDir + "/tpl"

	AppPatchDir = AppBaseDir + "/patch"

	AppDistDir = AppBaseDir + "/dist"

	AppTmpDir = AppBaseDir + "/tmp"
)

const (
	// DefaultClusterName is the default cluster name
	DefaultClusterName = "kubeon"

	DefaultClusterAPIDomain = "apiserver.cluster.local"

	DefaultClusterAPIPort = int32(6443)

	DefaultClusterLBMode = LocalHaproxy

	DefaultClusterDNSDomain = "cluster.local"

	DefaultClusterMaxPods = uint32(2000)

	DefaultClusterPortRange = "10000-30000"

	// DefaultSVCSubnet defines the default pod subnet
	DefaultSVCSubnet = "10.64.0.0/12"

	// DefaultPodSubnet defines the default pod subnet
	DefaultPodSubnet = "10.96.0.0/12"

	// DefaultImageRepo defines the default image repo
	DefaultImageRepo = "k8s.gcr.io"

	DefaultProxyMode = IPVSProxy

	DefaultIPVSScheduler = "rr"

	DefaultRuntimeMode = ContainerdRuntime

	DefaultNetworkMode = CalicoNetwork

	DefaultCalicoMode = CalicoIPIP

	DefaultCalicoMTU = "1440"

	KubeadmIgnorePreflightErrors = "SystemVerification,FileContent--proc-sys-net-bridge-bridge-nf-call-iptables"

	KubernetesEtcPath = "/etc/kubernetes"

	KubernetesPkiPath = "/etc/kubernetes/pki"

	KubernetesPkiEtcdPath = "/etc/kubernetes/pki/etcd"
)

var KubernetesDefaultConfigPath = onutil.K8sDir() + "/config"

const (
	UpdaterNamespace = "kube-system"

	HaproxyResource = "kubeon/local-haproxy"

	UpdaterResource = "kubeon/apiserver-updater"
)
