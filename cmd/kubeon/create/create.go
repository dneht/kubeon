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

package create

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"net"
	"os"
	"time"
)

type flagpole struct {
	define.DefaultList
	define.MasterList
	define.WorkerList
	MirrorHost                            string
	OnlyCreate                            bool
	UseOffline                            bool
	EnableIPv6                            bool
	ClusterLBDomain                       string
	ClusterDNSDomain                      string
	ClusterMaxPods                        uint32
	ClusterPortRange                      string
	ClusterNodeMaskSize                   uint32
	ClusterNodeMaskSizeV6                 uint32
	ClusterFeatureGates                   string
	ExternalLBIP                          string
	ExternalLBPort                        uint32
	InnerLBMode                           string
	NodeInterface                         []string
	NetworkSVCCIDR                        string
	NetworkSVCCIDRV6                      string
	NetworkPodCIDR                        string
	NetworkPodCIDRV6                      string
	InputProxyMode                        string
	IPVSScheduler                         string
	StrictARP                             bool
	InputCRIMode                          string
	InputCNIMode                          string
	InputICMode                           string
	WithNvidia                            bool
	WithKata                              bool
	WithKruise                            bool
	InputCertSANs                         []string
	CalicoMTU                             uint32
	CalicoEnableVXLAN                     bool
	CalicoEnableBPF                       bool
	CalicoEnableBGP                       bool
	CalicoEnableDSR                       bool
	CalicoEnableWireGuard                 bool
	CalicoEnableCrossSubnet               bool
	CalicoBPFBypassConntrack              bool
	CiliumMTU                             uint32
	CiliumEnableGENEVE                    bool
	CiliumEnableNR                        bool
	CiliumEnableBGP                       bool
	CiliumEnableDSR                       bool
	CiliumEnableBM                        bool
	CiliumEnableBBR                       bool
	CiliumEnableWireGuard                 bool
	CiliumNativeRoutingCIDR               string
	CiliumNativeRoutingCIDRV6             string
	CiliumEnableIPv6BigTCP                bool
	CiliumMonitorAggregation              string
	CiliumMonitorInterval                 string
	CiliumMonitorFlags                    string
	CiliumPolicyMode                      string
	CiliumEnableLocalRedirect             bool
	CiliumAutoProtectPortRange            bool
	CiliumLBNativeAcceleration            bool
	CiliumLBMaglevAlgorithm               bool
	CiliumEnableExternalClusterIP         bool
	CiliumBPFMapDynamicSizeRatio          string
	CiliumBPFLBMapMax                     uint32
	CiliumBPFPolicyMapMax                 uint32
	CiliumBPFHostNamespaceOnly            bool
	CiliumBPFBypassFIBLookup              bool
	CiliumInstallIptablesRules            bool
	CiliumInstallNoConntrackIptablesRules bool
	ContourDisableInsecure                bool
	ContourDisableMergeSlashes            bool
	IstioEnableAutoInject                 bool
	IstioServiceEntryExportTo             []string
	IstioVirtualServiceExportTo           []string
	IstioDestinationRuleExportTo          []string
	IstioEnableAutoMTLS                   bool
	IstioEnableHttp2AutoUpgrade           bool
	IstioEnablePrometheusMerge            bool
	IstioEnableNetworkPlugin              bool
	IstioEnableIngressGateway             bool
	IstioIngressGatewayType               string
	IstioEnableEgressGateway              bool
	IstioEgressGatewayType                string
	IstioEnableSkywalking                 bool
	IstioEnableSkywalkingAll              bool
	IstioSkywalkingService                string
	IstioSkywalkingPort                   uint32
	IstioSkywalkingAccessToken            string
	IstioEnableZipkin                     bool
	IstioZipkinService                    string
	IstioZipkinPort                       uint32
	IstioAccessLogServiceAddress          string
	IstioMetricsServiceAddress            string
	KruiseFeatureGates                    string
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "create CLUSTER_NAME CLUSTER_VERSION [flags]\n",
		Aliases: []string{"c", "new"},
		Short:   "Create a new cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return preRunE(flags, cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().StringVar(
		&flags.MirrorHost, "mirror",
		"yes", "default yes will use aliyun mirror, like: xxx.mirror.aliyuncs.com",
	)
	cmd.Flags().BoolVarP(
		&flags.OnlyCreate, "only-create", "C",
		false, "create config only",
	)
	cmd.Flags().BoolVarP(
		&flags.UseOffline, "offline", "O",
		false, "install use offline system package",
	)
	cmd.Flags().BoolVar(
		&flags.EnableIPv6, "enable-ipv6",
		false, "enable dual stack",
	)
	cmd.Flags().StringVar(
		&flags.ClusterLBDomain, "lb-domain",
		define.DefaultClusterAPIDomain,
		"A domain for apiserver load balance, not ip",
	)
	cmd.Flags().StringVar(
		&flags.ClusterDNSDomain, "dns-domain",
		define.DefaultClusterDNSDomain,
		"Cluster dns domain",
	)
	cmd.Flags().Uint32Var(
		&flags.ClusterMaxPods, "max-pods",
		define.DefaultClusterMaxPods,
		"Kubelet max pods",
	)
	cmd.Flags().StringVar(
		&flags.ClusterPortRange, "port-range",
		define.DefaultClusterPortRange,
		"Service node port range",
	)
	cmd.Flags().Uint32Var(
		&flags.ClusterNodeMaskSize, "node-mask-size",
		24,
		"controller manager node mask size",
	)
	cmd.Flags().Uint32Var(
		&flags.ClusterNodeMaskSizeV6, "node-mask-size-ipv6",
		64,
		"controller manager node mask size, ipv6",
	)
	cmd.Flags().StringVar(
		&flags.ClusterFeatureGates, "feature-gates",
		"",
		"controller manager feature gates",
	)
	cmd.Flags().StringSliceVar(
		&flags.NodeInterface, "interface",
		[]string{},
		"If the node has multiple network, specify one",
	)
	cmd.Flags().StringVar(
		&flags.NetworkSVCCIDR, "svc-cidr",
		define.DefaultSVCSubnet,
		"Use alternative range of IP address for service VIPs",
	)
	cmd.Flags().StringVar(
		&flags.NetworkSVCCIDRV6, "svc-cidr-ipv6",
		"",
		"Use alternative range of IP address for service VIPs, ipv6",
	)
	cmd.Flags().StringVar(
		&flags.NetworkPodCIDR, "pod-cidr",
		define.DefaultPodSubnet,
		"Specify range of IP addresses for the pod network",
	)
	cmd.Flags().StringVar(
		&flags.NetworkPodCIDRV6, "pod-cidr-ipv6",
		"",
		"Specify range of IP addresses for the pod network, ipv6",
	)
	cmd.Flags().StringVar(
		&flags.ExternalLBIP, "lb-ip",
		"",
		"External load balancer ip",
	)
	cmd.Flags().Uint32Var(
		&flags.ExternalLBPort, "lb-port",
		define.DefaultClusterAPIPort,
		"External load balancer port",
	)
	cmd.Flags().StringVar(
		&flags.InnerLBMode, "lb-mode",
		define.DefaultClusterLBMode,
		"Inner load balancer mode, only(haproxy or updater)",
	)
	cmd.Flags().StringVar(
		&flags.InputProxyMode, "proxy",
		define.DefaultProxyMode,
		"Kube proxy mode, only ipvs or iptables",
	)
	cmd.Flags().StringVar(
		&flags.IPVSScheduler, "ipvs-scheduler",
		define.DefaultIPVSScheduler,
		"IPVS scheduler, like: rr wrr",
	)
	cmd.Flags().BoolVar(
		&flags.StrictARP, "strict-arp",
		false,
		"IPVS strict ARP",
	)
	cmd.Flags().StringVar(
		&flags.InputCRIMode, "cri",
		define.DefaultRuntimeMode,
		"Runtime interface, only docker or containerd",
	)
	cmd.Flags().StringVar(
		&flags.InputCNIMode, "cni",
		define.DefaultNetworkMode,
		"Network interface, only none or calico",
	)
	cmd.Flags().StringVar(
		&flags.InputICMode, "ic",
		define.DefaultIngressMode,
		"Ingress controller, only none or contour",
	)
	cmd.Flags().BoolVar(
		&flags.WithNvidia, "with-nvidia",
		true,
		"Install nvidia",
	)
	cmd.Flags().BoolVar(
		&flags.WithKata, "with-kata",
		false,
		"Install kata with Kata-deploy",
	)
	cmd.Flags().BoolVar(
		&flags.WithKata, "with-kruise",
		false,
		"Install kruise",
	)
	cmd.Flags().StringSliceVar(
		&flags.InputCertSANs, "cert-san",
		[]string{},
		"Kubernetes api server CertSANs",
	)
	cmd.Flags().Uint32Var(
		&flags.CalicoMTU, "calico-mtu",
		0,
		"Calico MTU",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoEnableVXLAN, "calico-enable-vxlan",
		false,
		"Calico use vxlan, default if ipip",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoEnableBPF, "calico-enable-bpf",
		false,
		"Calico use bpf",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoEnableBGP, "calico-enable-bgp",
		false,
		"Calico use bgp",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoEnableDSR, "calico-enable-dsr",
		false,
		"Calico use dsr",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoEnableWireGuard, "calico-enable-wireguard",
		false,
		"Calico with wireguard",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoEnableCrossSubnet, "calico-cross-subnet",
		false,
		"If true, tunnel mode: Always -> CrossSubnet",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoBPFBypassConntrack, "calico-bpf-bypass-conntrack",
		true,
		"Calico BPF bypass conntrack",
	)
	cmd.Flags().Uint32Var(
		&flags.CiliumMTU, "cilium-mtu",
		0,
		"Cilium MTU",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableGENEVE, "cilium-enable-geneve",
		false,
		"Cilium use geneve, default if vxlan",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableNR, "cilium-enable-nr",
		false,
		"Cilium use native routing",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableBGP, "cilium-enable-bgp",
		false,
		"Cilium use bgp",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableDSR, "cilium-enable-dsr",
		false,
		"Cilium use dsr",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableBM, "cilium-enable-bm",
		false,
		"Cilium use bandwidth manager",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableBBR, "cilium-enable-bbr",
		false,
		"Cilium use bbr",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableWireGuard, "cilium-enable-wireguard",
		false,
		"Cilium with wireguard",
	)
	cmd.Flags().StringVar(
		&flags.CiliumNativeRoutingCIDR, "cilium-nr-cidr",
		"",
		"Native-Routing is required, unless bgp is enabled",
	)
	cmd.Flags().StringVar(
		&flags.CiliumNativeRoutingCIDRV6, "cilium-nr-cidr-ipv6",
		"",
		"Native-Routing and ipv6 is required, unless bgp is enabled",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableIPv6BigTCP, "cilium-enable-ipv6-bigtcp",
		false,
		"Cilium enable ipv6 big tcp",
	)
	cmd.Flags().StringVar(
		&flags.CiliumMonitorAggregation, "cilium-monitor-aggregation",
		"medium",
		"Cilium monitor aggregation, only: low, medium, maximum",
	)
	cmd.Flags().StringVar(
		&flags.CiliumMonitorInterval, "cilium-monitor-interval",
		"5s",
		"Only effective when monitor aggregation is set to medium or higher",
	)
	cmd.Flags().StringVar(
		&flags.CiliumMonitorFlags, "cilium-monitor-flags",
		"all",
		"Only effective when monitor aggregation is set to medium or higher",
	)
	cmd.Flags().StringVar(
		&flags.CiliumPolicyMode, "cilium-policy",
		"default",
		"Cilium policy mode, only: default, always, never",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableLocalRedirect, "cilium-enable-local-redirect",
		true,
		"Cilium enable local redirect",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumAutoProtectPortRange, "cilium-enable-protect-port-range",
		true,
		"Cilium auto protect nodePort range",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumLBNativeAcceleration, "cilium-lb-acc-native",
		false,
		"Cilium lb native acceleration",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumLBMaglevAlgorithm, "cilium-lb-use-maglev",
		false,
		"Cilium lb use maglev algorithm",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableExternalClusterIP, "cilium-enable-external-cluster-ip",
		true,
		"Cilium lb enable external clusterIP",
	)
	cmd.Flags().StringVar(
		&flags.CiliumBPFMapDynamicSizeRatio, "cilium-bpf-map-dynamic-size-ratio",
		"0.0025",
		"CiliumB bpf map dynamic size ratio",
	)
	cmd.Flags().Uint32Var(
		&flags.CiliumBPFLBMapMax, "cilium-bpf-lb-map-max",
		65536,
		"CiliumB bpf lb map max",
	)
	cmd.Flags().Uint32Var(
		&flags.CiliumBPFPolicyMapMax, "cilium-bpf-policy-map-max",
		16384,
		"CiliumB bpf policy map max",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumBPFHostNamespaceOnly, "cilium-bpf-hostns-only",
		false,
		"CiliumB bpf host namespace only",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumBPFBypassFIBLookup, "cilium-bpf-bypass-fib-lookup",
		false,
		"Cilium bpf bypass fib lookup",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumInstallIptablesRules, "cilium-install-iptables-rules",
		true,
		"Cilium install iptables rules",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumInstallNoConntrackIptablesRules, "cilium-install-no-conntrack-iptables-rules",
		true,
		"Cilium install no-conntrack iptables rules",
	)
	cmd.Flags().BoolVar(
		&flags.ContourDisableInsecure, "contour-disable-insecure",
		false,
		"Contour disable HTTPProxy permitInsecure",
	)
	cmd.Flags().BoolVar(
		&flags.ContourDisableMergeSlashes, "contour-disable-merge-slashes",
		false,
		"Contour disable envoy's non-standard merge_slashes path transformation option",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableNetworkPlugin, "istio-enable-cni",
		false,
		"Istio install cni plugin",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableAutoInject, "istio-enable-auto-inject",
		true,
		"Istio enable sidecar auto inject",
	)
	cmd.Flags().StringSliceVar(
		&flags.IstioServiceEntryExportTo, "istio-service-entry-export-to",
		[]string{},
		"Istio ServiceEntry export to namespace",
	)
	cmd.Flags().StringSliceVar(
		&flags.IstioVirtualServiceExportTo, "istio-virtual-service-export-to",
		[]string{},
		"Istio VirtualService export to namespace",
	)
	cmd.Flags().StringSliceVar(
		&flags.IstioDestinationRuleExportTo, "istio-destination-rule-export-to",
		[]string{},
		"Istio DestinationRule export to namespace",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableAutoMTLS, "istio-enable-auto-mtls",
		true,
		"Istio enable auto mTLS",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableHttp2AutoUpgrade, "istio-enable-http2-auto-upgrade",
		true,
		"Istio enable HTTP2 auto upgrade",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnablePrometheusMerge, "istio-enable-prometheus-merge",
		true,
		"Istio agent will merge metrics exposed by the application with metrics from Envoy and Istio agent",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableIngressGateway, "istio-enable-ingress",
		true,
		"Istio enable ingress gateway",
	)
	cmd.Flags().StringVar(
		&flags.IstioIngressGatewayType, "istio-ingress-gateway-type",
		"ClusterIP",
		"Istio ingress gateway load balance type",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableEgressGateway, "istio-enable-egress",
		false,
		"Istio enable egress gateway",
	)
	cmd.Flags().StringVar(
		&flags.IstioEgressGatewayType, "istio-egress-gateway-type",
		"ClusterIP",
		"Istio egress gateway load balance type",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableSkywalking, "istio-enable-skywalking",
		false,
		"Istio enable skywalking tracer",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableSkywalkingAll, "istio-enable-skywalking-all",
		false,
		"Istio enable skywalking tracer and metrics",
	)
	cmd.Flags().StringVar(
		&flags.IstioSkywalkingService, "istio-skywalking-svc",
		"",
		"Istio skywalking tracer service, not ip",
	)
	cmd.Flags().Uint32Var(
		&flags.IstioSkywalkingPort, "istio-skywalking-port",
		11800,
		"Istio skywalking tracer port",
	)
	cmd.Flags().StringVar(
		&flags.IstioSkywalkingAccessToken, "istio-skywalking-access-token",
		"",
		"Istio skywalking tracer access token",
	)
	cmd.Flags().BoolVar(
		&flags.IstioEnableZipkin, "istio-enable-zipkin",
		false,
		"Istio enable zipkin tracer",
	)
	cmd.Flags().StringVar(
		&flags.IstioSkywalkingService, "istio-zipkin-svc",
		"",
		"Istio zipkin tracer service, not ip",
	)
	cmd.Flags().Uint32Var(
		&flags.IstioSkywalkingPort, "istio-zipkin-port",
		9411,
		"Istio zipkin tracer port",
	)
	cmd.Flags().StringVar(
		&flags.IstioAccessLogServiceAddress, "istio-access-log-service-addr",
		"",
		"Istio accessLog service address",
	)
	cmd.Flags().StringVar(
		&flags.IstioMetricsServiceAddress, "istio-metrics-service-addr",
		"",
		"Istio metrics service address",
	)
	cmd.Flags().StringVar(
		&flags.KruiseFeatureGates, "kruise-feature-gates",
		"",
		"Kruise feature gates, like: g1,g2,g3...",
	)
	cmd.Flags().UintVar(
		&flags.DefaultPort, "default-port",
		22, "All node default port",
	)
	cmd.Flags().StringVar(
		&flags.DefaultUser, "default-user",
		"root", "All node default username",
	)
	cmd.Flags().StringVar(
		&flags.DefaultPassword, "default-passwd",
		"", "All node default password",
	)
	cmd.Flags().StringVar(
		&flags.DefaultPkFile, "default-pk-file",
		"", "All node default private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringVar(
		&flags.DefaultPkPassword, "default-pk-password",
		"", "All node default private key file password",
	)
	cmd.Flags().IPSliceVarP(
		&flags.MasterIPs, "master", "m",
		[]net.IP{}, "Multi master ips",
	)
	cmd.Flags().UintSliceVar(
		&flags.MasterPorts, "master-port",
		[]uint{}, "Multi master port, default is 22",
	)
	cmd.Flags().StringSliceVarP(
		&flags.MasterNames, "master-name", "M",
		[]string{}, "Multi master name, if using will go to set hostname, otherwise use hostname",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterLabels, "master-labels",
		[]string{}, "Multi master labels",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterUsers, "master-user",
		[]string{}, "Multi master username",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterPasswords, "master-passwd",
		[]string{}, "Multi master password",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterPkFiles, "master-pk-file",
		[]string{}, "Multi master private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterPkPasswords, "master-pk-password",
		[]string{}, "Multi master private key file password",
	)
	cmd.Flags().UintVar(
		&flags.MasterDefaultPort, "default-master-port",
		22, "Multi master default port",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultUser, "default-master-user",
		"root", "Multi master default username",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultPassword, "default-master-passwd",
		"", "Multi master default password",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultPkFile, "default-master-pk-file",
		"", "Multi master default private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultPkPassword, "default-master-pk-password",
		"", "Multi master default private key file password",
	)
	cmd.Flags().IPSliceVarP(
		&flags.WorkerIPs, "worker", "w",
		[]net.IP{}, "Multi worker ips",
	)
	cmd.Flags().UintSliceVar(
		&flags.WorkerPorts, "worker-port",
		[]uint{}, "Multi worker port, default is 22",
	)
	cmd.Flags().StringSliceVarP(
		&flags.WorkerNames, "worker-name", "W",
		[]string{}, "Multi worker name, if using will go to set hostname, otherwise use hostname",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterLabels, "worker-labels",
		[]string{}, "Multi worker labels",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerUsers, "worker-user",
		[]string{}, "Multi worker username",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerPasswords, "worker-passwd",
		[]string{}, "Multi worker password",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerPkFiles, "worker-pk-file",
		[]string{}, "Multi master private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerPkPasswords, "worker-pk-password",
		[]string{}, "Multi worker private key file password",
	)
	cmd.Flags().UintVar(
		&flags.WorkerDefaultPort, "default-worker-port",
		22, "Multi worker default port",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultUser, "default-worker-user",
		"root", "Multi worker default username",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultPassword, "default-worker-passwd",
		"", "Multi worker default password",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultPkFile, "default-worker-pk-file",
		"", "Multi worker default private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultPkPassword, "default-worker-pk-password",
		"", "Multi worker default private key file password",
	)
	return cmd
}

func preRunE(flags *flagpole, cmd *cobra.Command, args []string) error {
	inputVersion, err := define.NewStdVersion(args[1])
	if nil != err {
		return err
	}
	version := inputVersion.Full
	if !checkSupport(flags, version) ||
		!flags.MasterList.CheckMatch() || !flags.WorkerList.CheckMatch() {
		os.Exit(1)
	}
	proxyReplace := preFlags(flags)
	klog.V(1).Infof("Create cluster with cri=%s, cni=%s, proxy=%s, ingress=%s, nvidia=%v, kata=%v, kruise=%v",
		flags.InputCRIMode, flags.InputCNIMode, flags.InputProxyMode, flags.InputICMode,
		flags.WithNvidia, flags.WithKata, flags.WithKruise)

	current := &cluster.Cluster{
		Version:        inputVersion,
		Mirror:         onutil.ConvMirror(flags.MirrorHost, define.MirrorImageRepo, define.DockerImageRepo),
		IsOffline:      flags.UseOffline,
		EnableDual:     flags.EnableIPv6,
		LbPort:         flags.ExternalLBPort,
		LbDomain:       flags.ClusterLBDomain,
		LbMode:         flags.InnerLBMode,
		DnsDomain:      flags.ClusterDNSDomain,
		MaxPods:        flags.ClusterMaxPods,
		PortRange:      flags.ClusterPortRange,
		NodeMaskSize:   flags.ClusterNodeMaskSize,
		NodeMaskSizeV6: flags.ClusterNodeMaskSizeV6,
		FeatureGates:   flags.ClusterFeatureGates,
		SvcCIDR:        flags.NetworkSVCCIDR,
		SvcV6CIDR:      flags.NetworkSVCCIDRV6,
		PodCIDR:        flags.NetworkPodCIDR,
		PodV6CIDR:      flags.NetworkPodCIDRV6,
		NodeInterface:  flags.NodeInterface,
		ProxyMode:      flags.InputProxyMode,
		IPVSScheduler:  flags.IPVSScheduler,
		StrictARP:      flags.StrictARP,
		RuntimeMode:    flags.InputCRIMode,
		NetworkMode:    flags.InputCNIMode,
		EnableBPF:      proxyReplace,
		CalicoConf: &cluster.CalicoConf{
			CalicoMTU:          flags.CalicoMTU,
			EnableVXLAN:        flags.CalicoEnableVXLAN,
			EnableBPF:          flags.CalicoEnableBPF,
			EnableDSR:          flags.CalicoEnableDSR,
			EnableBGP:          flags.CalicoEnableBGP,
			EnableWireGuard:    flags.CalicoEnableWireGuard,
			EnableCrossSubnet:  flags.CalicoEnableCrossSubnet,
			BPFBypassConntrack: flags.CalicoBPFBypassConntrack,
		},
		CiliumConf: &cluster.CiliumConf{
			CiliumMTU:                       flags.CiliumMTU,
			EnableGENEVE:                    flags.CiliumEnableGENEVE,
			EnableNR:                        flags.CiliumEnableNR,
			EnableDSR:                       flags.CiliumEnableDSR,
			EnableBGP:                       flags.CiliumEnableBGP,
			EnableBM:                        flags.CiliumEnableBM,
			EnableBBR:                       flags.CiliumEnableBBR,
			EnableWireGuard:                 flags.CiliumEnableWireGuard,
			NativeRoutingCIDR:               flags.CiliumNativeRoutingCIDR,
			NativeRoutingCIDRV6:             flags.CiliumNativeRoutingCIDRV6,
			EnableIPv6BigTCP:                flags.CiliumEnableIPv6BigTCP,
			MonitorAggregation:              flags.CiliumMonitorAggregation,
			MonitorFlags:                    flags.CiliumMonitorFlags,
			MonitorInterval:                 flags.CiliumMonitorInterval,
			PolicyMode:                      flags.CiliumPolicyMode,
			MapDynamicSizeRatio:             flags.CiliumBPFMapDynamicSizeRatio,
			PolicyMapMax:                    flags.CiliumBPFPolicyMapMax,
			LBMapMax:                        flags.CiliumBPFLBMapMax,
			EnableLocalRedirect:             flags.CiliumEnableLocalRedirect,
			AutoProtectPortRange:            flags.CiliumAutoProtectPortRange,
			LBNativeAcceleration:            flags.CiliumLBNativeAcceleration,
			LBMaglevAlgorithm:               flags.CiliumLBMaglevAlgorithm,
			EnableExternalClusterIP:         flags.CiliumEnableExternalClusterIP,
			BPFHostNamespaceOnly:            flags.CiliumBPFHostNamespaceOnly,
			BPFBypassFIBLookup:              flags.CiliumBPFBypassFIBLookup,
			InstallIptablesRules:            flags.CiliumInstallIptablesRules,
			InstallNoConntrackIptablesRules: flags.CiliumInstallNoConntrackIptablesRules,
		},
		IngressMode: flags.InputICMode,
		ContourConf: &cluster.ContourConf{
			DisableInsecure:     flags.ContourDisableInsecure,
			DisableMergeSlashes: flags.ContourDisableMergeSlashes,
		},
		IstioConf: &cluster.IstioConf{
			EnableNetworkPlugin:     flags.IstioEnableNetworkPlugin,
			EnableAutoInject:        flags.IstioEnableAutoInject,
			ServiceEntryExportTo:    flags.IstioServiceEntryExportTo,
			VirtualServiceExportTo:  flags.IstioVirtualServiceExportTo,
			DestinationRuleExportTo: flags.IstioDestinationRuleExportTo,
			EnableAutoMTLS:          flags.IstioEnableAutoMTLS,
			EnableHttp2AutoUpgrade:  flags.IstioEnableHttp2AutoUpgrade,
			EnablePrometheusMerge:   flags.IstioEnablePrometheusMerge,
			EnableIngressGateway:    flags.IstioEnableIngressGateway,
			IngressGatewayType:      flags.IstioIngressGatewayType,
			EnableEgressGateway:     flags.IstioEnableEgressGateway,
			EgressGatewayType:       flags.IstioEgressGatewayType,
			EnableSkywalking:        flags.IstioEnableSkywalking,
			EnableSkywalkingAll:     flags.IstioEnableSkywalkingAll,
			SkywalkingService:       flags.IstioSkywalkingService,
			SkywalkingPort:          flags.IstioSkywalkingPort,
			SkywalkingAccessToken:   flags.IstioSkywalkingAccessToken,
			EnableZipkin:            flags.IstioEnableZipkin,
			ZipkinService:           flags.IstioZipkinService,
			ZipkinPort:              flags.IstioZipkinPort,
			AccessLogServiceAddress: flags.IstioAccessLogServiceAddress,
			MetricsServiceAddress:   flags.IstioMetricsServiceAddress,
		},
		UseNvidia: flags.WithNvidia,
		UseKata:   flags.WithKata,
		UseKruise: flags.WithKruise,
		KruiseConf: &cluster.KruiseConf{
			FeatureGates: flags.KruiseFeatureGates,
		},
		CertSANs: flags.InputCertSANs,
		Status:   cluster.StatusCreating,
	}
	klog.V(1).Info("Ready to check & prepare host, please wait a moment...")
	return cluster.InitNewCluster(current, flags.ExternalLBIP, flags.DefaultList, flags.MasterList, flags.WorkerList)
}

func preFlags(flags *flagpole) bool {
	if "" == flags.NetworkSVCCIDRV6 || "" == flags.NetworkPodCIDRV6 {
		flags.EnableIPv6 = false
	}
	flags.WithNvidia = flags.WithNvidia && flags.InputCRIMode == define.ContainerdRuntime

	if flags.CiliumEnableGENEVE || flags.CiliumEnableNR || flags.CiliumEnableDSR || flags.CiliumEnableBM || flags.CiliumEnableBBR {
		flags.InputCNIMode = define.CiliumNetwork
	}
	if flags.InputICMode == define.IstioIngress {
		flags.CiliumBPFHostNamespaceOnly = true
	}
	if flags.CiliumEnableNR {
		flags.CiliumEnableGENEVE = false
	}
	if flags.CiliumEnableDSR {
		flags.CiliumEnableGENEVE = false
		flags.CiliumEnableNR = true
	} else {
		flags.CiliumInstallNoConntrackIptablesRules = false
	}

	if flags.CalicoEnableVXLAN || flags.CalicoEnableBPF || flags.CalicoEnableBGP || flags.CalicoEnableDSR {
		flags.InputCNIMode = define.CalicoNetwork
	}
	if flags.CalicoEnableBGP {
		flags.CalicoEnableVXLAN = false
	}
	if flags.CalicoEnableDSR {
		flags.CalicoEnableVXLAN = false
		flags.CalicoEnableBPF = true
	}

	if "" == flags.IstioSkywalkingService || flags.IstioSkywalkingPort == 0 {
		flags.IstioEnableSkywalking = false
		flags.IstioEnableSkywalkingAll = false
	}
	if "" == flags.IstioZipkinService || flags.IstioZipkinPort == 0 {
		flags.IstioEnableZipkin = false
	}

	if flags.CiliumEnableNR && "" == flags.CiliumNativeRoutingCIDR {
		klog.Error("If native routing is used, the cilium-nr-cidr must be set")
		os.Exit(1)
	}
	proxyReplace := flags.CalicoEnableBPF || flags.InputCNIMode == define.CiliumNetwork
	if proxyReplace {
		flags.InputProxyMode = flags.InputCNIMode
	}
	return proxyReplace
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) (err error) {
	current := cluster.Current()
	if nil == current {
		return errors.New("cluster create error")
	}
	if flags.OnlyCreate {
		return nil
	}

	err = preInstall(current, current.Mirror)
	if nil != err {
		klog.Warningf("Prepare install failed, please check: %v", err)
		return nil
	}
	err = initCluster(current)
	if nil != err {
		klog.Errorf("Create cluster failed, reset nodes: %v", err)
		action.KubeadmResetList(cluster.CurrentNodes(), false, false)
	}
	return nil
}

func preInstall(current *cluster.Cluster, mirror string) (err error) {
	err = cluster.CreateResource(mirror)
	if nil != err {
		return err
	}

	err = module.PrepareInstall(cluster.CurrentNodes(), false)
	if nil != err {
		return err
	}
	return nil
}

func initCluster(current *cluster.Cluster) (err error) {
	bootNode := cluster.BootstrapNode()
	err = module.SetupBootKubeadm(bootNode)
	if nil != err {
		return err
	}
	err = module.InstallInner(define.HealthzReader)
	if nil != err {
		return err
	}
	err = module.InstallNetwork()
	if nil != err {
		return err
	}
	err = module.InstallExtend()
	if nil != err {
		return err
	}
	err = action.KubeadmInitWait(4 * time.Minute)
	if nil != err {
		return err
	}

	joinNodes := cluster.JoinsNode()
	err = module.SetupJoinsKubeadm(joinNodes)
	if nil != err {
		return err
	}
	err = module.InstallLoadBalance(current.Workers)
	if nil != err {
		return err
	}
	module.LabelDevice()
	err = module.InstallIngress()
	if nil != err {
		klog.V(1).Info("Cluster has been installed, but the ingress installation failed, please use `kubeon view conf` to generate a configuration file to retry")
	}
	return cluster.CreateCompleteCluster()
}
