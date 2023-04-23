/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package module

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"strings"
)

func InstallExtend(isUpgrade bool) (err error) {
	current := cluster.Current()
	if current.UseNvidia && current.HasNvidia {
		err = InstallInner(define.NvidiaRuntime, isUpgrade)
		if nil != err {
			return err
		}
	}
	if current.UseKata {
		err = InstallInner(define.KataRuntime, isUpgrade)
		if nil != err {
			return err
		}
	}
	if current.UseKruise {
		err = InstallInner(define.KruisePlugin, isUpgrade)
		if nil != err {
			return err
		}
	}
	return nil
}

func InstallNetwork(isUpgrade bool) (err error) {
	current := cluster.Current()
	switch current.NetworkMode {
	case define.CalicoNetwork:
		{
			err = InstallInner(define.CalicoNetwork, isUpgrade)
			if nil != err {
				return err
			}
		}
	case define.CiliumNetwork:
		{
			err = InstallInner(define.CiliumNetwork, isUpgrade)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

func InstallIngress(isUpgrade bool) (err error) {
	current := cluster.Current()
	switch current.IngressMode {
	case define.ContourIngress:
		{
			err = InstallInner(define.ContourIngress, isUpgrade)
			if nil != err {
				return err
			}
			break
		}
	case define.IstioIngress:
		{
			err = InstallInner(define.IstioIngress, isUpgrade)
			if nil != err {
				return err
			}
			break
		}
	}
	return nil
}

func InstallInner(moduleName string, isUpgrade bool) (err error) {
	current := cluster.Current()
	local := current.IsRealLocal()
	switch moduleName {
	case define.CiliumNetwork:
		ciliumTpl := ciliumTemplate(current)
		if isUpgrade {
			err = action.CiliumExecute(release.BuildCiliumUpgradeArgs(ciliumTpl, local))
		} else {
			err = action.CiliumExecute(release.BuildCiliumInstallArgs(ciliumTpl, local))
			if nil != err {
				return err
			}
			err = action.CiliumExecute(release.BuildHubbleInstallArgs(ciliumTpl, local))
		}
	case define.IstioIngress:
		istioTpl := istioTemplate(current)
		err = action.KubectlCreateNamespace(define.IstioNamespace)
		if nil != err {
			return err
		}
		err = action.IstioExecute(release.BuildIstioInstallArgs(istioTpl, local))
	default:
		var bytes []byte
		bytes, err = ShowInner(moduleName)
		if nil != err {
			return err
		}
		if nil != bytes {
			klog.V(4).Infof("Install %s on cluster", moduleName)
			err = action.KubectlApplyData(bytes)
		}
	}
	return err
}

func DeleteInner(moduleName string) (err error) {
	bytes, err := ShowInner(moduleName)
	if nil != err {
		return err
	}
	if nil != bytes {
		return action.KubectlDeleteData(bytes)
	}
	return nil
}

func getKubeletTemplate() *release.KubeletTemplate {
	current := cluster.Current()
	return &release.KubeletTemplate{
		APIVersion:       current.GetKubeletAPIVersion(),
		ClusterDnsIP:     current.DnsIP,
		ClusterDnsDomain: current.DnsDomain,
		ClusterMaxPods:   current.MaxPods,
	}
}

func ShowInner(moduleName string) ([]byte, error) {
	current := cluster.Current()
	local := current.IsRealLocal()
	klog.V(4).Infof("[module] Get module [%s] config", moduleName)
	switch moduleName {
	case define.KubeletModule:
		return release.RenderKubeletTemplate(getKubeletTemplate(), current.Version.Full)
	case define.CorednsPart:
		return release.RenderCorednsTemplate(&release.CorednsTemplate{
			CPVersion:    current.GetModuleVersion(define.CorednsPart),
			MirrorUrl:    current.Mirror,
			ClusterDnsIP: current.DnsIP,
		}, local)
	case define.CalicoNetwork:
		if nil == current.CalicoConf {
			return nil, errors.New("get calico config error")
		}
		return release.RenderCalicoTemplate(calicoTemplate(current), local)
	case define.CiliumNetwork:
		if nil == current.CiliumConf {
			return nil, errors.New("get cilium config error")
		}
		return release.RenderCiliumCommand(ciliumTemplate(current), local)
	case define.ContourIngress:
		if nil == current.ContourConf {
			return nil, errors.New("get contour config error")
		}
		return release.RenderContourTemplate(contourTemplate(current), local)
	case define.IstioIngress:
		if nil == current.IstioConf {
			return nil, errors.New("get istio config error")
		}
		return release.RenderIstioTemplate(istioTemplate(current), local)
	case define.NvidiaRuntime:
		return release.RenderNvidiaTemplate(&release.NvidiaTemplate{
			CPVersion: current.GetModuleVersion(define.NvidiaRuntime),
			MirrorUrl: current.Mirror,
		}, local)
	case define.KataRuntime:
		return release.RenderKataTemplate(&release.KataTemplate{
			CPVersion: current.GetModuleVersion(define.KataRuntime),
			MirrorUrl: current.Mirror,
		}, local)
	case define.KruisePlugin:
		return release.RenderKruiseTemplate(&release.KruiseTemplate{
			CPVersion: current.GetModuleVersion(define.KruisePlugin),
			MirrorUrl: current.Mirror,
		}, local)
	case define.HealthzReader:
		return release.RenderHealthzTemplate(current.Version.Full), nil
	case define.LocalHaproxy:
		return release.RenderHaproxyTemplate(&release.BalanceHaproxyTemplate{
			MasterHosts: current.MasterAPIs(),
			ImageUrl:    current.GetHaproxyImageAddr(),
		})
	case define.ApiserverStartup:
		return release.RenderStartupService(&release.ApiserverScriptTemplate{
			TargetDomain: current.LbDomain,
			VirtualAddr:  current.LbIP,
			RealAddrs:    strings.Join(current.MasterIPs(), ","),
		})
	case define.ApiserverUpdater:
		return release.RenderUpdaterTemplate(&release.ApiserverUpdaterTemplate{
			ClusterLbIP: current.LbIP,
			MasterIPs:   current.MasterIPs(),
			ImageUrl:    current.GetUpdaterImageAddr(),
		})
	default:
		klog.Warningf("Not support inner module[%s]", moduleName)
		return nil, nil
	}
}

func getInterface(interfaces []string) (bool, string) {
	isSetInterface := len(interfaces) != 0
	defaultInterface := ""
	if isSetInterface {
		defaultInterface = interfaces[0]
	}
	return isSetInterface, defaultInterface
}

func calicoTemplate(current *cluster.Cluster) *release.CalicoTemplate {
	calicoConf := current.CalicoConf
	isSetInterface, defaultInterface := getInterface(current.NodeInterface)
	backendMode, ipipMode, vxlanMode, vxlanv6Mode := calicoMode(current, calicoConf)
	lbMode := calicoLBMode(calicoConf)
	return &release.CalicoTemplate{
		CPVersion:             current.GetModuleVersion(define.CalicoNetwork),
		MirrorUrl:             current.Mirror,
		IsSetInterface:        isSetInterface,
		DefaultInterface:      defaultInterface,
		BackendMode:           backendMode,
		EnableBPF:             calicoConf.EnableBPF,
		EnableWireGuard:       calicoConf.EnableWireGuard,
		CalicoMTU:             calicoConf.CalicoMTU,
		LBMode:                lbMode,
		EnableVXLAN:           calicoConf.EnableVXLAN,
		IPIPMode:              ipipMode,
		VXLANMode:             vxlanMode,
		VXLANv6Mode:           vxlanv6Mode,
		EnablePassConntrack:   calicoConf.EnablePassConntrack,
		ClusterEnableDual:     current.EnableDual,
		ClusterLbDomain:       current.LbDomain,
		ClusterLbPort:         current.LbPort,
		ClusterPortRange:      current.PortRange,
		ClusterNodeMaskSize:   current.NodeMaskSize,
		ClusterNodeMaskSizeV6: current.NodeMaskSizeV6,
		ClusterPodCIDR:        current.PodCIDR,
		ClusterPodCIDRV6:      current.PodV6CIDR,
	}
}

func calicoMode(current *cluster.Cluster, conf *cluster.CalicoConf) (string, string, string, string) {
	backendMode := define.CalicoBackendBIRD
	ipipMode := define.CalicoTunModeNever
	vxlanMode := define.CalicoTunModeNever
	vxlanv6Mode := define.CalicoTunModeNever

	if conf.EnableBGP {
		return backendMode, ipipMode, vxlanMode, vxlanv6Mode
	}
	if conf.EnableVXLAN {
		backendMode = define.CalicoBackendVXLAN
		if conf.EnableCrossSubnet {
			vxlanMode = define.CalicoTunModeCrossSubnet
		} else {
			vxlanMode = define.CalicoTunModeAlways
		}
		if current.EnableDual {
			vxlanv6Mode = vxlanMode
		}
	} else {
		if conf.EnableCrossSubnet {
			ipipMode = define.CalicoTunModeCrossSubnet
		} else {
			ipipMode = define.CalicoTunModeAlways
		}
	}
	return backendMode, ipipMode, vxlanMode, vxlanv6Mode
}

func calicoLBMode(conf *cluster.CalicoConf) string {
	lbMode := define.CalicoLBModeDefault
	if conf.EnableDSR {
		lbMode = define.CalicoLBModeDSR
	}
	return lbMode
}

func ciliumTemplate(current *cluster.Cluster) *release.CiliumTemplate {
	ciliumConf := current.CiliumConf
	isSetInterface, defaultInterface := getInterface(current.NodeInterface)
	ipv4Masq, ipv4NRCIDR, ipv6Masq, ipv6NRCIDR, tunnelMode, autoDirectNode := ciliumMode(current, ciliumConf)
	lbMode := ciliumLBMode(ciliumConf)
	return &release.CiliumTemplate{
		CPVersion:               current.GetModuleVersion(define.CiliumNetwork),
		MirrorUrl:               current.Mirror,
		IsSetInterface:          isSetInterface,
		DefaultInterface:        defaultInterface,
		EnableBGP:               ciliumConf.EnableBGP,
		EnableBM:                ciliumConf.EnableBM,
		EnableBBR:               ciliumConf.EnableBBR,
		EnableWireGuard:         ciliumConf.EnableWireGuard,
		EnableIPv4Masquerade:    ipv4Masq,
		EnableIPv6Masquerade:    ipv6Masq,
		NativeRoutingCIDR:       ipv4NRCIDR,
		NativeRoutingCIDRV6:     ipv6NRCIDR,
		CiliumMTU:               ciliumConf.CiliumMTU,
		PolicyMode:              ciliumConf.PolicyMode,
		TunnelMode:              tunnelMode,
		LBMode:                  lbMode,
		EnableEndpointRoutes:    ciliumConf.EnableEndpointRoutes,
		EnableLocalRedirect:     ciliumConf.EnableLocalRedirect,
		EnableHostnsOnly:        ciliumConf.EnableHostnsOnly,
		AutoDirectNodeRoutes:    autoDirectNode,
		EnableEndpointSlice:     ciliumConf.EnableEndpointSlice,
		EnableExternalClusterIP: ciliumConf.EnableExternalClusterIP,
		AutoProtectPortRange:    ciliumConf.AutoProtectPortRange,
		HubbleVersion:           current.GetModuleVersion(define.CiliumHubble),
		CustomConfigs:           ciliumConf.CustomConfigs,
		ClusterEnableDual:       current.EnableDual,
		ClusterLbDomain:         current.LbDomain,
		ClusterLbPort:           current.LbPort,
		ClusterPortRange:        strings.ReplaceAll(current.PortRange, "-", ","),
		ClusterNodeMaskSize:     current.NodeMaskSize,
		ClusterNodeMaskSizeV6:   current.NodeMaskSizeV6,
		ClusterPodCIDR:          current.PodCIDR,
		ClusterPodCIDRV6:        current.PodV6CIDR,
	}
}

func ciliumMode(current *cluster.Cluster, conf *cluster.CiliumConf) (bool, string, bool, string, string, bool) {
	ipv4NRCIDR := ""
	ipv6NRCIDR := ""
	tunnelMode := define.CiliumTunnelVXLAN
	autoDirectNode := false

	if conf.EnableBGP {
		tunnelMode = define.CiliumTunnelDisabled
		autoDirectNode = true
		return false, ipv4NRCIDR, false, ipv6NRCIDR, tunnelMode, autoDirectNode
	}
	if conf.EnableGENEVE {
		tunnelMode = define.CiliumTunnelGENEVE
	}
	if conf.EnableNR {
		if "" == conf.NativeRoutingCIDR {
			ipv4NRCIDR = current.PodCIDR
		} else {
			ipv4NRCIDR = conf.NativeRoutingCIDR
		}
		if current.EnableDual {
			if "" == conf.NativeRoutingCIDR {
				ipv6NRCIDR = current.PodV6CIDR
			} else {
				ipv6NRCIDR = conf.NativeRoutingCIDRV6
			}
		}
		tunnelMode = define.CiliumTunnelDisabled
		autoDirectNode = true
	}
	if !conf.AutoDirectNodeRoutes {
		autoDirectNode = false
	}
	return conf.EnableIPv4Masquerade, ipv4NRCIDR, conf.EnableIPv6Masquerade, ipv6NRCIDR, tunnelMode, autoDirectNode
}

func ciliumLBMode(conf *cluster.CiliumConf) string {
	lbMode := define.CiliumLBModeSNAT
	if conf.EnableNR {
		lbMode = define.CiliumLBModeHybrid
	}
	if conf.EnableDSR {
		lbMode = define.CiliumLBModeDSR
	}
	return lbMode
}

func contourTemplate(current *cluster.Cluster) *release.ContourTemplate {
	contourConf := current.ContourConf
	return &release.ContourTemplate{
		CPVersion:             current.GetModuleVersion(define.ContourIngress),
		MirrorUrl:             current.Mirror,
		DisableMergeSlashes:   contourConf.DisableMergeSlashes,
		DisablePermitInsecure: contourConf.DisableInsecure,
	}
}

func istioTemplate(current *cluster.Cluster) *release.IstioTemplate {
	istioConf := current.IstioConf
	return &release.IstioTemplate{
		CPVersion:               current.GetModuleVersion(define.IstioIngress),
		MirrorUrl:               current.Mirror,
		EnableAutoInject:        istioConf.EnableAutoInject,
		ServiceEntryExportTo:    istioConf.ServiceEntryExportTo,
		VirtualServiceExportTo:  istioConf.VirtualServiceExportTo,
		DestinationRuleExportTo: istioConf.DestinationRuleExportTo,
		IngressGatewayType:      istioConf.IngressGatewayType,
		EnableEgressGateway:     istioConf.EnableEgressGateway,
		EgressGatewayType:       istioConf.EgressGatewayType,
		CustomConfigs:           istioConf.CustomConfigs,
		ProxyClusterDomain:      current.DnsDomain,
	}
}
