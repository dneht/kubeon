/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package release

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/define"
	"strconv"
	"strings"
)

func BuildCiliumInstallArgs(input *CiliumTemplate, local bool, size int) []string {
	return buildCiliumInstallArgs(input, true, local, size)
}

func BuildCiliumUpgradeArgs(input *CiliumTemplate, local bool, size int) []string {
	return buildCiliumUpgradeArgs(input, true, local, size)
}

func BuildHubbleInstallArgs(input *CiliumTemplate, local bool) []string {
	return buildHubbleInstallArgs(input, true, local)
}

func BuildIstioInstallArgs(input *IstioTemplate, local bool) []string {
	return buildIstioInstallArgs(input, true, local)
}

func buildCiliumInstallArgs(input *CiliumTemplate, install, local bool, size int) []string {
	imageHub := input.MirrorUrl + define.CiliumMirrorImagePrefix
	ciliumOperatorImage := define.CiliumOperatorMirrorImage
	if local {
		imageHub = define.QuayImageRepo + define.CiliumImagePrefix
		ciliumOperatorImage = define.CiliumOperatorImage
	}
	ciliumArgs := make([]string, 0, 8)
	if install {
		ciliumArgs = append(ciliumArgs, "install")
	} else {
		ciliumArgs = append(ciliumArgs, "install")
	}
	waitTime := 30 * int64(size)
	if waitTime < 300 {
		waitTime = 300
	}
	ciliumArgs = append(ciliumArgs, "--version", input.CPVersion,
		"--wait-duration", strconv.FormatInt(waitTime, 10)+"s",
		"--agent-image", imageHub+"/"+define.CiliumAgentImage+":"+input.CPVersion,
		"--operator-image", imageHub+"/"+ciliumOperatorImage+":"+input.CPVersion)
	if input.EnableWireGuard {
		ciliumArgs = append(ciliumArgs, "--encryption", "wireguard")
	}
	ciliumArgs = append(ciliumArgs, "--helm-set", "kubeProxyReplacement=strict",
		"--helm-set", "k8sServiceHost="+input.ClusterLbDomain,
		"--helm-set", "k8sServicePort="+strconv.FormatUint(uint64(input.ClusterLbPort), 10),
		"--helm-set", "ipam.mode=kubernetes",
		"--helm-set", "wellKnownIdentities.enabled=true",
		"--helm-set", "tunnel="+input.TunnelMode,
		"--helm-set", "hostPort.enabled=true", "--helm-set", "nodePort.enabled=true",
		"--helm-set", "hubble.enabled=true", "--helm-set", "hubble.relay.enabled=true", "--helm-set", "hubble.ui.enabled=true",
		"--helm-set", "prometheus.enabled=true", "--helm-set", "operator.prometheus.enabled=true",
	)
	ciliumConfs := make([]string, 0, 32)
	for _, conf := range input.CustomConfigs {
		ciliumConfs = append(ciliumConfs, conf)
	}
	ciliumConfs = append(ciliumConfs, "allow-localhost=always",
		"sockops-enable=false", "enable-k8s-endpoint-slice=true", "enable-k8s-event-handover=true",
		"enable-external-ips=true", "enable-host-port=true", "enable-node-port=true")
	if input.IsSetInterface {
		ciliumConfs = append(ciliumConfs, "devices="+input.DefaultInterface)
	}
	if input.EnableBGP {
		ciliumConfs = append(ciliumConfs, "enable-bgp-control-plane=true")
	}
	if input.EnableBM {
		ciliumConfs = append(ciliumConfs, "enable-bandwidth-manager=true")
	}
	if input.EnableBBR {
		ciliumConfs = append(ciliumConfs, "enable-bbr=true")
	}
	ciliumConfs = append(ciliumConfs, "enable-ipv4=true")
	if "" != input.ClusterPodCIDR {
		ciliumConfs = append(ciliumConfs, "cluster-pool-ipv4-cidr="+input.ClusterPodCIDR)
	}
	if input.ClusterNodeMaskSize > 0 {
		ciliumConfs = append(ciliumConfs, "cluster-pool-ipv4-mask-size="+strconv.FormatUint(uint64(input.ClusterNodeMaskSize), 10))
	}
	if input.EnableIPv4Masquerade {
		ciliumConfs = append(ciliumConfs, "enable-ipv4-masquerade=true")
		if "" != input.NativeRoutingCIDR {
			ciliumConfs = append(ciliumConfs, "ipv4-native-routing-cidr="+input.NativeRoutingCIDR)
		} else {
			ciliumConfs = append(ciliumConfs, "ipv4-native-routing-cidr="+input.ClusterPodCIDR)
		}
	} else {
		ciliumConfs = append(ciliumConfs, "enable-ipv4-masquerade=false")
	}
	if input.ClusterEnableDual {
		ciliumConfs = append(ciliumConfs, "enable-ipv6=true")
		if "" != input.ClusterPodCIDRV6 {
			ciliumConfs = append(ciliumConfs, "cluster-pool-ipv6-cidr="+input.ClusterPodCIDRV6)
		}
		if input.ClusterNodeMaskSizeV6 > 0 {
			ciliumConfs = append(ciliumConfs, "cluster-pool-ipv6-mask-size="+strconv.FormatUint(uint64(input.ClusterNodeMaskSizeV6), 10))
		}
		if input.EnableIPv6Masquerade {
			ciliumConfs = append(ciliumConfs, "enable-ipv6-masquerade=true")
			if "" != input.NativeRoutingCIDRV6 {
				ciliumConfs = append(ciliumConfs, "ipv6-native-routing-cidr="+input.NativeRoutingCIDRV6)
			} else {
				ciliumConfs = append(ciliumConfs, "ipv6-native-routing-cidr="+input.ClusterPodCIDRV6)
			}
		} else {
			ciliumConfs = append(ciliumConfs, "enable-ipv6-masquerade=false")
		}
	} else {
		ciliumConfs = append(ciliumConfs, "enable-ipv6=false", "enable-ipv6-masquerade=false")
	}
	ciliumConfs = append(ciliumConfs, "enable-bpf-masquerade=true")
	if input.CiliumMTU > 0 {
		ciliumConfs = append(ciliumConfs, "mtu="+strconv.FormatUint(uint64(input.CiliumMTU), 10))
	}
	ciliumConfs = append(ciliumConfs, "enable-policy="+input.PolicyMode, "tunnel="+input.TunnelMode)
	if input.EnableEndpointRoutes {
		ciliumConfs = append(ciliumConfs, "enable-endpoint-routes=true")
	} else {
		ciliumConfs = append(ciliumConfs, "enable-endpoint-routes=false")
	}
	if input.EnableLocalRedirect {
		ciliumConfs = append(ciliumConfs, "enable-local-redirect-policy=true")
	} else {
		ciliumConfs = append(ciliumConfs, "enable-local-redirect-policy=false")
	}
	ciliumConfs = append(ciliumConfs, "bpf-lb-mode="+input.LBMode, "bpf-lb-sock=true")
	if input.EnableHostnsOnly {
		ciliumConfs = append(ciliumConfs, "bpf-lb-sock-hostns-only=true")
	} else {
		ciliumConfs = append(ciliumConfs, "bpf-lb-sock-hostns-only=false")
	}
	if input.EnableExternalClusterIP {
		ciliumConfs = append(ciliumConfs, "bpf-lb-external-clusterip=true")
	} else {
		ciliumConfs = append(ciliumConfs, "bpf-lb-external-clusterip=false")
	}
	if input.AutoDirectNodeRoutes {
		ciliumConfs = append(ciliumConfs, "auto-direct-node-routes=true")
	} else {
		ciliumConfs = append(ciliumConfs, "auto-direct-node-routes=false")
	}
	if input.EnableEndpointSlice {
		ciliumConfs = append(ciliumConfs, "enable-cilium-endpoint-slice=true")
	} else {
		ciliumConfs = append(ciliumConfs, "enable-cilium-endpoint-slice=false")
	}
	if input.AutoProtectPortRange {
		ciliumConfs = append(ciliumConfs, "enable-auto-protect-node-port-range=true")
	} else {
		ciliumConfs = append(ciliumConfs, "enable-auto-protect-node-port-range=false")
	}
	ciliumArgs = append(ciliumArgs, "--config", strings.Join(ciliumConfs, ","))
	return ciliumArgs
}

func buildCiliumUpgradeArgs(input *CiliumTemplate, install, local bool, size int) []string {
	imageHub := input.MirrorUrl + define.CiliumMirrorImagePrefix
	ciliumOperatorImage := define.CiliumOperatorMirrorImage
	if local {
		imageHub = define.QuayImageRepo + define.CiliumImagePrefix
		ciliumOperatorImage = define.CiliumOperatorImage
	}
	ciliumArgs := make([]string, 0, 8)
	if install {
		ciliumArgs = append(ciliumArgs, "upgrade")
	} else {
		ciliumArgs = append(ciliumArgs, "upgrade")
	}
	waitTime := 30 * int64(size)
	if waitTime < 300 {
		waitTime = 300
	}
	ciliumArgs = append(ciliumArgs, "--version", input.CPVersion,
		"--wait-duration", strconv.FormatInt(waitTime, 10)+"s",
		"--agent-image", imageHub+"/"+define.CiliumAgentImage+":"+input.CPVersion,
		"--operator-image", imageHub+"/"+ciliumOperatorImage+":"+input.CPVersion,
		"--hubble-relay-image", imageHub+"/"+define.HubbleRelayImage+":"+input.CPVersion)
	return ciliumArgs
}

func buildHubbleInstallArgs(input *CiliumTemplate, install, local bool) []string {
	imageHub := input.MirrorUrl + define.CiliumMirrorImagePrefix
	if local {
		imageHub = define.QuayImageRepo + define.CiliumImagePrefix
	}
	ciliumArgs := make([]string, 0, 8)
	if install {
		ciliumArgs = append(ciliumArgs, "hubble", "enable")
	} else {
		ciliumArgs = append(ciliumArgs, "hubble", "enable")
	}
	ciliumArgs = append(ciliumArgs, "--relay", "--relay-version", input.CPVersion,
		"--relay-image", imageHub+"/"+define.HubbleRelayImage+":"+input.CPVersion,
		"--ui", "--ui-version", input.HubbleVersion,
		"--ui-image", imageHub+"/"+define.HubbleUIImage+":"+input.HubbleVersion,
		"--ui-backend-image", imageHub+"/"+define.HubbleUIBackendImage+":"+input.HubbleVersion)
	ciliumArgs = append(ciliumArgs,
		"--helm-set", "k8sServiceHost="+input.ClusterLbDomain,
		"--helm-set", "k8sServicePort="+strconv.FormatUint(uint64(input.ClusterLbPort), 10),
		"--helm-set", "tunnel="+input.TunnelMode,
		"--helm-set", "hubble.relay.enabled=true", "--helm-set", "hubble.ui.enabled=true",
		"--helm-set", "hubble.metrics.serviceMonitor.enabled=true")
	return ciliumArgs
}

func buildIstioInstallArgs(input *IstioTemplate, install, local bool) []string {
	imageHub := input.MirrorUrl + define.IstioMirrorImagePrefix
	imagePullPolicy := define.ImagePullPolicyAlways
	istioProxyImage := define.IstioProxyImage
	istioPilotImage := define.IstioPilotImage
	istioInstallCNIImage := define.IstioCNIImage
	if local {
		imageHub = define.DockerImageRepo + define.IstioImagePrefix
		imagePullPolicy = define.ImagePullPolicyNotPresent
		istioProxyImage = define.IstioMirrorProxyImage
		istioPilotImage = define.IstioMirrorPilotImage
		istioInstallCNIImage = define.IstioMirrorCNIImage
	}
	istioArgs := make([]string, 0, 64)
	if install {
		istioArgs = append(istioArgs, "install")
	} else {
		istioArgs = append(istioArgs, "manifest", "generate")
	}
	for _, conf := range input.CustomConfigs {
		istioArgs = append(istioArgs, "--set", conf)
	}

	autoInject := define.IstioProxyAutoInjectEnable
	if !input.EnableAutoInject {
		autoInject = define.IstioProxyAutoInjectDisable
	}
	clusterDomain := input.ProxyClusterDomain
	if "" == clusterDomain {
		clusterDomain = define.DefaultClusterDNSDomain
	}
	istioArgs = append(istioArgs, "--set", "profile=default", "--set", "values.global.hub="+imageHub,
		"--set", "values.global.proxy.autoInject="+autoInject, "--set", "values.global.proxy.clusterDomain="+clusterDomain,
		"--set", "values.global.imagePullPolicy="+imagePullPolicy,
		"--set", "values.global.proxy.image="+istioProxyImage, "--set", "values.global.proxy_init.image="+istioProxyImage,
		"--set", "values.pilot.image="+istioPilotImage, "--set", "values.cni.image="+istioInstallCNIImage,
	)
	//if input.EnableNetworkPlugin {
	//	istioArgs = append(istioArgs, "--set", "components.cni.enabled=true",
	//		"--set", "values.cni.excludeNamespaces[0]=kube-system", "--set", "values.cni.excludeNamespaces[1]=istio-system")
	//}
	//istioArgs = append(istioArgs, "--set", "meshConfig.enableAutoMtls="+strconv.FormatBool(input.EnableAutoMTLS))
	//if input.EnableHttp2AutoUpgrade {
	//	istioArgs = append(istioArgs, "--set", "meshConfig.h2UpgradePolicy="+define.IstioHttp2AutoUpgrade)
	//} else {
	//	istioArgs = append(istioArgs, "--set", "meshConfig.h2UpgradePolicy="+define.IstioHttp2DontAutoUpgrade)
	//}
	istioArgs = append(istioArgs, "--set", "values.gateways.istio-ingressgateway.enabled=true",
		"--set", "values.gateways.istio-ingressgateway.type="+input.IngressGatewayType)
	if input.EnableEgressGateway {
		istioArgs = append(istioArgs, "--set", "values.gateways.istio-egressgateway.enabled=true")
		if "" != input.EgressGatewayType {
			istioArgs = append(istioArgs, "--set", "values.gateways.istio-egressgateway.type="+input.EgressGatewayType)
		}
	} else {
		istioArgs = append(istioArgs, "--set", "values.gateways.istio-egressgateway.enabled=false")
	}
	if len(input.ServiceEntryExportTo) > 0 {
		for idx, exportTo := range input.ServiceEntryExportTo {
			istioArgs = append(istioArgs, "--set", fmt.Sprintf("meshConfig.defaultServiceExportTo[%d]=%s", idx, exportTo))
		}
	}
	if len(input.VirtualServiceExportTo) > 0 {
		for idx, exportTo := range input.VirtualServiceExportTo {
			istioArgs = append(istioArgs, "--set", fmt.Sprintf("meshConfig.defaultVirtualServiceExportTo[%d]=%s", idx, exportTo))
		}
	}
	if len(input.DestinationRuleExportTo) > 0 {
		for idx, exportTo := range input.DestinationRuleExportTo {
			istioArgs = append(istioArgs, "--set", fmt.Sprintf("meshConfig.defaultDestinationRuleExportTo[%d]=%s", idx, exportTo))
		}
	}
	istioArgs = append(istioArgs, "--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[0]=.*membership_healthy.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[1]=.*upstream_cx_active.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[2]=.*upstream_cx_total.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[3]=.*upstream_rq_active.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[4]=.*upstream_rq_total.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[5]=.*upstream_rq_pending_active.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[6]=.*lb_healthy_panic.*",
		"--set", "meshConfig.defaultConfig.proxyStatsMatcher.inclusionRegexps[7]=.*upstream_cx_none_healthy.*")
	return istioArgs
}
