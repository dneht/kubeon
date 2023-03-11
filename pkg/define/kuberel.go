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
	K8S_1_25_0, _ = NewStdVersion("v1.25.0")
	K8S_1_19_x, _ = NewRngVersion("v1.19.16", "v1.19.16")
	K8S_1_20_x, _ = NewRngVersion("v1.20.15", "v1.20.15")
	K8S_1_21_x, _ = NewRngVersion("v1.21.14", "v1.21.14")
	K8S_1_22_x, _ = NewRngVersion("v1.22.15", "v1.22.17")
	K8S_1_23_x, _ = NewRngVersion("v1.23.11", "v1.23.17")
	K8S_1_24_x, _ = NewRngVersion("v1.24.1", "v1.24.11")
	K8S_1_25_x, _ = NewRngVersion("v1.25.1", "v1.25.7")
	K8S_1_26_x, _ = NewRngVersion("v1.26.1", "v1.26.2")
	ETCD_3_4_0, _ = NewStdVersion("3.4.0")
)

const (
	KubeletConfigApiB1 = "v1beta1"
	KubeadmConfigApiB2 = "v1beta2"
	KubeadmConfigApiB3 = "v1beta3"
)

const (
	DefaultImageRepo = "registry.k8s.io"
	DockerImageRepo  = "docker.io"
	QuayImageRepo    = "quay.io"
	MirrorImageRepo  = "uhub.service.ucloud.cn"
	AliyunImageRepo  = "registry.cn-hangzhou.aliyuncs.com"
)

const (
	ImagePullPolicyAlways     = "Always"
	ImagePullPolicyNotPresent = "IfNotPresent"
	ImagePullPolicyNever      = "Never"
)

const (
	ServiceClusterIP    = "ClusterIP"
	ServiceNodePort     = "NodePort"
	ServiceLoadBalancer = "LoadBalancer"
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

const (
	HealthzReader    = "healthz"
	LocalHaproxy     = "haproxy"
	ApiserverUpdater = "updater"
	ApiserverStartup = "startup"
	ApiserverService = "apiserver-startup"
	BPFMountService  = "sys-bpf-mount"
)

const InstallScript = "script"

const (
	DockerRuntime     = "docker"
	ContainerdRuntime = "containerd"
	NvidiaRuntime     = "nvidia"
	KataRuntime       = "kata"
)

const (
	CorednsPart = "coredns"
)

const (
	IPTablesProxy = "iptables"
	IPVSProxy     = "ipvs"
	BPFProxy      = "bpf"
)

const (
	NetworkPlugin = "cni"

	NoneNetwork   = "none"
	CalicoNetwork = "calico"
	CiliumNetwork = "cilium"
	CiliumHubble  = "hubble"

	NoneIngress    = "none"
	IstioIngress   = "istio"
	ContourIngress = "contour"
)

const (
	KruisePlugin = "kruise"
)

const (
	CalicoBackendBIRD  = "bird"
	CalicoBackendVXLAN = "vxlan"

	CalicoTunnelIPIP  = "ipip"
	CalicoTunnelVXLAN = "vxlan"

	CalicoLBModeDSR     = "DSR"
	CalicoLBModeDefault = "Tunnel"

	CalicoTunModeNever       = "Never"
	CalicoTunModeAlways      = "Always"
	CalicoTunModeCrossSubnet = "CrossSubnet"
)

const (
	CiliumCommand = "/opt/cni/bin/cilium"

	CiliumImagePrefix         = "/cilium"
	CiliumMirrorImagePrefix   = "/kubeon"
	CiliumAgentImage          = "cilium"
	CiliumOperatorImage       = "operator-generic"
	CiliumOperatorMirrorImage = "cilium-operator-generic"
	HubbleRelayImage          = "hubble-relay"
	HubbleUIImage             = "hubble-ui"
	HubbleUIBackendImage      = "hubble-ui-backend"

	CiliumPolicyDefault = "default"
	CiliumPolicyAlways  = "always"
	CiliumPolicyNever   = "never"

	CiliumTunnelDisabled = "disabled"
	CiliumTunnelVXLAN    = "vxlan"
	CiliumTunnelGENEVE   = "geneve"

	CiliumLBModeDSR    = "dsr"
	CiliumLBModeSNAT   = "snat"
	CiliumLBModeHybrid = "hybrid"

	CiliumLBAlgorithmMaglev = "maglev"
	CiliumLBAlgorithmRandom = "random"

	CiliumLBAccelerationDisabled = "disabled"
	CiliumLBAccelerationNative   = "native"
)

const (
	ContourNamespace = "projectcontour"
)

const (
	IstioNamespace = "istio-system"

	IstioCommand = "/opt/cni/bin/istioctl"

	IstioImagePrefix       = "/istio"
	IstioMirrorImagePrefix = "/kubeon"
	IstioProxyImage        = "proxyv2"
	IstioMirrorProxyImage  = "istio-proxyv2"
	IstioPilotImage        = "pilot"
	IstioMirrorPilotImage  = "istio-pilot"
	IstioCNIImage          = "install-cni"
	IstioMirrorCNIImage    = "istio-install-cni"

	IstioProxyAutoInjectEnable  = "enabled"
	IstioProxyAutoInjectDisable = "disabled"

	IstioHttp2AutoUpgrade     = "UPGRADE"
	IstioHttp2DontAutoUpgrade = "DO_NOT_UPGRADE"
)
