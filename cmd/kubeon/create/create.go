/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package create

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cloud"
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
	CalicoEnablePassConntrack             bool
	CiliumMTU                             uint32
	CiliumEnableGENEVE                    bool
	CiliumEnableNR                        bool
	CiliumEnableBGP                       bool
	CiliumEnableDSR                       bool
	CiliumEnableBM                        bool
	CiliumEnableBBR                       bool
	CiliumEnableWireGuard                 bool
	CiliumDisableIpMasq                   bool
	CiliumDisableIpMasqV6                 bool
	CiliumNativeRoutingCIDR               string
	CiliumNativeRoutingCIDRV6             string
	CiliumDisableEndpointRoutes           bool
	CiliumDisableLocalRedirect            bool
	CiliumDisableHostnsOnly               bool
	CiliumDisableAutoDirectNodeRoutes     bool
	CiliumDisableEndpointSlice            bool
	CiliumDisableExternalClusterIP        bool
	CiliumDisableAutoProtectNodePortRange bool
	CiliumPolicyMode                      string
	CiliumConfig                          string
	ContourDisableInsecure                bool
	ContourDisableMergeSlashes            bool
	IstioDisableAutoInject                bool
	IstioServiceEntryExportTo             []string
	IstioVirtualServiceExportTo           []string
	IstioDestinationRuleExportTo          []string
	IstioIngressGatewayType               string
	IstioEnableEgressGateway              bool
	IstioEgressGatewayType                string
	IstioConfig                           string
	NvidiaElevated                        bool
	NvidiaDriverRoot                      string
	KruiseFeatureGates                    string
	CloudProvider                         string
	CloudEndpoint                         string
	CloudRouterTableIds                   []string
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
		"yes", "default yes will use cn mirror",
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
		false,
		"Install nvidia",
	)
	cmd.Flags().BoolVar(
		&flags.WithKata, "with-kata",
		false,
		"Install kata with Kata-deploy",
	)
	cmd.Flags().BoolVar(
		&flags.WithKruise, "with-kruise",
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
		&flags.CalicoEnableCrossSubnet, "calico-enable-cross-subnet",
		false,
		"If true, tunnel mode: Always -> CrossSubnet",
	)
	cmd.Flags().BoolVar(
		&flags.CalicoEnablePassConntrack, "calico-enable-pass-conntrack",
		false,
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
		"Cilium use bbr, if kernel >= 5.18",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumEnableWireGuard, "cilium-enable-wireguard",
		false,
		"Cilium with wireguard",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableIpMasq, "cilium-disable-ip-masq",
		false,
		"Cilium enable ipv4 masquerade",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableIpMasqV6, "cilium-disable-ipv6-masq",
		false,
		"Cilium enable ipv6 masquerade",
	)
	cmd.Flags().StringVar(
		&flags.CiliumNativeRoutingCIDR, "cilium-nr-cidr",
		"",
		"Native-Routing is required, unless bgp is enabled",
	)
	cmd.Flags().StringVar(
		&flags.CiliumNativeRoutingCIDRV6, "cilium-nr-ipv6-cidr",
		"",
		"Native-Routing and ipv6 is required, unless bgp is enabled",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableEndpointRoutes, "cilium-disable-endpoint-routes",
		false,
		"Cilium disable endpoint routes",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableLocalRedirect, "cilium-disable-local-redirect",
		false,
		"Cilium disable local redirect",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableHostnsOnly, "cilium-disable-hostns-only",
		false,
		"Cilium disable hostns only",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableAutoDirectNodeRoutes, "cilium-disable-auto-direct-node-routes",
		false,
		"Cilium disable auto direct node routes",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableEndpointSlice, "cilium-disable-endpoint-slice",
		false,
		"If true, cilium will not create CiliumEndpoint",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableExternalClusterIP, "cilium-disable-external-clusterip",
		false,
		"If true, cilium will deny external access to cluster ip ",
	)
	cmd.Flags().BoolVar(
		&flags.CiliumDisableAutoProtectNodePortRange, "cilium-disable-auto-protect-node-port-range",
		false,
		"If true, cilium will not add the NodePort range ports to the kernel parameters to avoid the occupation of other applications",
	)
	cmd.Flags().StringVar(
		&flags.CiliumPolicyMode, "cilium-policy",
		"default",
		"Cilium policy mode, only: default, always, never",
	)
	cmd.Flags().StringVar(
		&flags.CiliumConfig, "cilium-config",
		"",
		"Set ConfigMap entries { key=value[,key=value,..] }",
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
		&flags.IstioDisableAutoInject, "istio-disable-auto-inject",
		false,
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
	cmd.Flags().StringVar(
		&flags.IstioConfig, "istio-config",
		"",
		"Set IstioOperator value { key=value[,key=value,..] }",
	)
	cmd.Flags().BoolVar(
		&flags.NvidiaElevated, "nvidia-elevated",
		false,
		"Deploy the nvidia daemonset with elevated privileges",
	)
	cmd.Flags().StringVar(
		&flags.NvidiaDriverRoot, "nvidia-driver-root",
		"/",
		"Nvidia driver root",
	)
	cmd.Flags().StringVar(
		&flags.KruiseFeatureGates, "kruise-feature-gates",
		"",
		"Kruise feature gates, like: g1,g2,g3...",
	)
	cmd.Flags().StringVar(
		&flags.CloudProvider, "cloud-provider",
		"",
		"Cloud provider",
	)
	cmd.Flags().StringVar(
		&flags.CloudEndpoint, "cloud-endpoint",
		"",
		"Cloud endpoint",
	)
	cmd.Flags().StringSliceVar(
		&flags.CloudRouterTableIds, "cloud-router-table-id",
		[]string{},
		"Cloud router table id",
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
	proxyReplace := preFlags(flags, inputVersion)
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
			CalicoMTU:           flags.CalicoMTU,
			EnableVXLAN:         flags.CalicoEnableVXLAN,
			EnableBPF:           flags.CalicoEnableBPF,
			EnableDSR:           flags.CalicoEnableDSR,
			EnableBGP:           flags.CalicoEnableBGP,
			EnableWireGuard:     flags.CalicoEnableWireGuard,
			EnableCrossSubnet:   flags.CalicoEnableCrossSubnet,
			EnablePassConntrack: flags.CalicoEnablePassConntrack,
		},
		CiliumConf: &cluster.CiliumConf{
			CiliumMTU:               flags.CiliumMTU,
			EnableGENEVE:            flags.CiliumEnableGENEVE,
			EnableNR:                flags.CiliumEnableNR,
			EnableDSR:               flags.CiliumEnableDSR,
			EnableBGP:               flags.CiliumEnableBGP,
			EnableBM:                flags.CiliumEnableBM,
			EnableBBR:               flags.CiliumEnableBBR,
			EnableWireGuard:         flags.CiliumEnableWireGuard,
			EnableIPv4Masquerade:    !flags.CiliumDisableIpMasq,
			EnableIPv6Masquerade:    !flags.CiliumDisableIpMasqV6,
			NativeRoutingCIDR:       flags.CiliumNativeRoutingCIDR,
			NativeRoutingCIDRV6:     flags.CiliumNativeRoutingCIDRV6,
			EnableEndpointRoutes:    !flags.CiliumDisableEndpointRoutes,
			EnableLocalRedirect:     !flags.CiliumDisableLocalRedirect,
			EnableHostnsOnly:        !flags.CiliumDisableHostnsOnly,
			AutoDirectNodeRoutes:    !flags.CiliumDisableAutoDirectNodeRoutes,
			EnableEndpointSlice:     !flags.CiliumDisableEndpointSlice,
			EnableExternalClusterIP: !flags.CiliumDisableExternalClusterIP,
			AutoProtectPortRange:    !flags.CiliumDisableAutoProtectNodePortRange,
			PolicyMode:              flags.CiliumPolicyMode,
			CustomConfigs:           checkConfigs(flags.CiliumConfig),
		},
		IngressMode: flags.InputICMode,
		ContourConf: &cluster.ContourConf{
			DisableInsecure:     flags.ContourDisableInsecure,
			DisableMergeSlashes: flags.ContourDisableMergeSlashes,
		},
		IstioConf: &cluster.IstioConf{
			EnableAutoInject:        !flags.IstioDisableAutoInject,
			ServiceEntryExportTo:    flags.IstioServiceEntryExportTo,
			VirtualServiceExportTo:  flags.IstioVirtualServiceExportTo,
			DestinationRuleExportTo: flags.IstioDestinationRuleExportTo,
			IngressGatewayType:      flags.IstioIngressGatewayType,
			EnableEgressGateway:     flags.IstioEnableEgressGateway,
			EgressGatewayType:       flags.IstioEgressGatewayType,
			CustomConfigs:           checkConfigs(flags.IstioConfig),
		},
		UseNvidia: flags.WithNvidia,
		NvidiaConf: &cluster.NvidiaConf{
			Elevated:   flags.NvidiaElevated,
			DriverRoot: flags.NvidiaDriverRoot,
		},
		UseKata:   flags.WithKata,
		UseKruise: flags.WithKruise,
		KruiseConf: &cluster.KruiseConf{
			FeatureGates: flags.KruiseFeatureGates,
		},
		CloudProvider: flags.CloudProvider,
		CloudConf: &cluster.CloudConf{
			Endpoint:       flags.CloudEndpoint,
			RouterTableIds: flags.CloudRouterTableIds,
		},
		CertSANs: flags.InputCertSANs,
		Status:   cluster.StatusCreating,
	}
	klog.V(1).Info("Ready to check & prepare host, please wait a moment...")
	return cluster.InitNewCluster(current, flags.ExternalLBIP, flags.DefaultList, flags.MasterList, flags.WorkerList)
}

func preFlags(flags *flagpole, version *define.StdVersion) bool {
	if "" == flags.NetworkSVCCIDRV6 || "" == flags.NetworkPodCIDRV6 {
		flags.EnableIPv6 = false
	}
	flags.WithNvidia = (flags.WithNvidia || flags.NvidiaElevated) && flags.InputCRIMode == define.ContainerdRuntime && version.IsSupportNvidia()
	flags.WithKata = flags.WithKata && version.IsSupportKata()

	proxyReplace := false
	if flags.InputCNIMode == define.BPFProxy {
		proxyReplace = true
		flags.InputCNIMode = define.NoneNetwork
		flags.InputProxyMode = define.BPFProxy
	}
	if flags.InputCNIMode != define.NoneNetwork {
		if flags.CiliumEnableGENEVE || flags.CiliumEnableNR || flags.CiliumEnableDSR || flags.CiliumEnableBM || flags.CiliumEnableBBR {
			flags.InputCNIMode = define.CiliumNetwork
		}
		if flags.InputICMode == define.IstioIngress {
			flags.CiliumDisableHostnsOnly = false
		}
		if flags.CiliumEnableNR {
			flags.CiliumEnableGENEVE = false
		}
		if flags.CiliumEnableDSR {
			flags.CiliumEnableGENEVE = false
			flags.CiliumEnableNR = true
		}
		if flags.CiliumEnableBBR {
			flags.CiliumEnableBM = true
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

		proxyReplace = flags.CalicoEnableBPF || flags.InputCNIMode == define.CiliumNetwork
		if proxyReplace {
			flags.InputProxyMode = flags.InputCNIMode
		}
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
	err = module.InstallInner(define.HealthzReader, false)
	if nil != err {
		return err
	}
	klog.V(1).Info("Preparing to initialize the network, please wait a moment...")
	err = module.InstallNetwork(false)
	if nil != err {
		return err
	}
	err = module.InstallDevice(false)
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
	if current.IsOnCloud() {
		cloud.ModifyRouterNow()
	}
	err = module.InstallLoadBalance(current.Workers)
	if nil != err {
		return err
	}
	module.LabelDevice()
	err = module.InstallIngress(false)
	if nil != err {
		klog.V(1).Info("Cluster has been installed, but the ingress installation failed, please use `kubeon view conf` to generate a configuration file to retry")
	}
	err = module.InstallExtend(false)
	if nil != err {
		return err
	}
	return cluster.CreateCompleteCluster()
}
