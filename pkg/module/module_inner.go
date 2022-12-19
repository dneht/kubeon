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

func InstallExtend() (err error) {
	current := cluster.Current()
	if current.UseNvidia && current.HasNvidia {
		err = InstallInner(define.NvidiaRuntime)
		if nil != err {
			return err
		}
	}
	if current.UseKata {
		err = InstallInner(define.KataRuntime)
		if nil != err {
			return err
		}
	}
	return nil
}

func InstallNetwork() (err error) {
	current := cluster.Current()
	switch current.NetworkMode {
	case define.CalicoNetwork:
		{
			err = InstallInner(define.CalicoNetwork)
			if nil != err {
				return err
			}
		}
	case define.CiliumNetwork:
		{
			err = InstallInner(define.CiliumNetwork)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

func InstallIngress() (err error) {
	current := cluster.Current()
	switch current.IngressMode {
	case define.ContourIngress:
		{
			err = InstallInner(define.ContourIngress)
			if nil != err {
				return err
			}
			break
		}
	case define.IstioIngress:
		{
			err = action.KubectlCreateNamespace(define.IstioNamespace)
			if nil != err {
				return err
			}
			err = InstallInner(define.IstioIngress)
			if nil != err {
				return err
			}
			break
		}
	}
	return nil
}

func InstallInner(moduleName string) (err error) {
	bytes, err := ShowInner(moduleName)
	if nil != err {
		return err
	}
	if nil != bytes {
		klog.V(4).Infof("Install %s on cluster", moduleName)
		return action.KubectlApplyData(bytes)
	}
	return nil
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
			MirrorUrl:    current.Mirror,
			ClusterDnsIP: current.DnsIP,
		}, local)
	case define.CalicoNetwork:
		if nil == current.CalicoConf {
			return nil, errors.New("get calico config error")
		}
		calicoConf := current.CalicoConf
		isSetInterface, defaultInterface := getInterface(current.NodeInterface)
		backendMode, ipipMode, vxlanMode, vxlanv6Mode := calicoMode(calicoConf, current.EnableDual)
		lbMode := calicoLBMode(calicoConf)
		return release.RenderCalicoTemplate(&release.CalicoTemplate{
			MirrorUrl:              current.Mirror,
			IsSetInterface:         isSetInterface,
			DefaultInterface:       defaultInterface,
			BackendMode:            backendMode,
			EnableBPF:              calicoConf.EnableBPF,
			EnableWireGuard:        calicoConf.EnableWireGuard,
			CalicoMTU:              calicoConf.CalicoMTU,
			LBMode:                 lbMode,
			EnableVXLAN:            calicoConf.EnableVXLAN,
			IPIPMode:               ipipMode,
			VXLANMode:              vxlanMode,
			VXLANv6Mode:            vxlanv6Mode,
			BPFHostConntrackBypass: calicoConf.BPFBypassConntrack,
			ClusterEnableDual:      current.EnableDual,
			ClusterLbDomain:        current.LbDomain,
			ClusterLbPort:          current.LbPort,
			ClusterPortRange:       current.PortRange,
			ClusterNodeMaskSize:    current.NodeMaskSize,
			ClusterNodeMaskSizeV6:  current.NodeMaskSizeV6,
			ClusterPodCIDR:         current.PodCIDR,
			ClusterPodCIDRV6:       current.PodV6CIDR,
		}, local)
	case define.CiliumNetwork:
		if nil == current.CiliumConf {
			return nil, errors.New("get cilium config error")
		}
		ciliumConf := current.CiliumConf
		isSetInterface, defaultInterface := getInterface(current.NodeInterface)
		ipv4Masq, ipv4NRCIDR, ipv6Masq, ipv6NRCIDR, tunnelMode, autoDirectNode := ciliumMode(ciliumConf, current.EnableDual)
		lbMode, lbAcceleration, lbAlgorithm := ciliumLBMode(ciliumConf)
		return release.RenderCiliumTemplate(&release.CiliumTemplate{
			MirrorUrl:                       current.Mirror,
			IsSetInterface:                  isSetInterface,
			DefaultInterface:                defaultInterface,
			EnableBGP:                       ciliumConf.EnableBGP,
			EnableBM:                        ciliumConf.EnableBM,
			EnableBBR:                       ciliumConf.EnableBBR,
			EnableWireGuard:                 ciliumConf.EnableWireGuard,
			EnableIPv4Masquerade:            ipv4Masq,
			EnableIPv6Masquerade:            ipv6Masq,
			NativeRoutingCIDR:               ipv4NRCIDR,
			NativeRoutingCIDRV6:             ipv6NRCIDR,
			EnableIPv6BigTCP:                ciliumConf.EnableIPv6BigTCP,
			CiliumMTU:                       ciliumConf.CiliumMTU,
			TunnelMode:                      tunnelMode,
			PolicyMode:                      ciliumConf.PolicyMode,
			LBMode:                          lbMode,
			LBAcceleration:                  lbAcceleration,
			LBAlgorithm:                     lbAlgorithm,
			LBHostNamespaceOnly:             ciliumConf.BPFHostNamespaceOnly,
			AutoDirectNodeRoutes:            autoDirectNode,
			EnableLocalRedirect:             ciliumConf.EnableLocalRedirect,
			AutoProtectPortRange:            ciliumConf.AutoProtectPortRange,
			BPFMapDynamicSizeRatio:          ciliumConf.MapDynamicSizeRatio,
			BPFLBMapMax:                     ciliumConf.LBMapMax,
			BPFPolicyMapMax:                 ciliumConf.PolicyMapMax,
			BPFLBExternalClusterIP:          ciliumConf.EnableExternalClusterIP,
			BPFLBBypassFIBLookup:            ciliumConf.BPFBypassFIBLookup,
			InstallIptablesRules:            ciliumConf.InstallIptablesRules,
			InstallNoConntrackIptablesRules: ciliumConf.InstallNoConntrackIptablesRules,
			MonitorAggregation:              ciliumConf.MonitorAggregation,
			MonitorInterval:                 ciliumConf.MonitorInterval,
			MonitorFlags:                    ciliumConf.MonitorFlags,
			ClusterEnableDual:               current.EnableDual,
			ClusterLbDomain:                 current.LbDomain,
			ClusterLbPort:                   current.LbPort,
			ClusterPortRange:                strings.ReplaceAll(current.PortRange, "-", ","),
			ClusterNodeMaskSize:             current.NodeMaskSize,
			ClusterNodeMaskSizeV6:           current.NodeMaskSizeV6,
			ClusterPodCIDR:                  current.PodCIDR,
			ClusterPodCIDRV6:                current.PodV6CIDR,
		}, local)
	case define.ContourIngress:
		if nil == current.ContourConf {
			return nil, errors.New("get contour config error")
		}
		contourConf := current.ContourConf
		return release.RenderContourTemplate(&release.ContourTemplate{
			MirrorUrl:             current.Mirror,
			DisableMergeSlashes:   contourConf.DisableMergeSlashes,
			DisablePermitInsecure: contourConf.DisableInsecure,
		}, local)
	case define.IstioIngress:
		if nil == current.IstioConf {
			return nil, errors.New("get contour config error")
		}
		istioConf := current.IstioConf
		return release.RenderIstioTemplate(&release.IstioTemplate{
			MirrorUrl:               current.Mirror,
			ProxyClusterDomain:      current.DnsDomain,
			EnableAutoInject:        istioConf.EnableAutoInject,
			ServiceEntryExportTo:    istioConf.ServiceEntryExportTo,
			VirtualServiceExportTo:  istioConf.VirtualServiceExportTo,
			DestinationRuleExportTo: istioConf.DestinationRuleExportTo,
			EnableAutoMTLS:          istioConf.EnableAutoMTLS,
			EnableHttp2AutoUpgrade:  istioConf.EnableHttp2AutoUpgrade,
			EnablePrometheusMerge:   istioConf.EnablePrometheusMerge,
			EnableNetworkPlugin:     istioConf.EnableNetworkPlugin,
			EnableIngressGateway:    istioConf.EnableIngressGateway,
			IngressGatewayType:      istioConf.IngressGatewayType,
			EnableEgressGateway:     istioConf.EnableEgressGateway,
			EgressGatewayType:       istioConf.EgressGatewayType,
			EnableSkywalking:        istioConf.EnableSkywalking,
			EnableSkywalkingAll:     istioConf.EnableSkywalkingAll,
			SkywalkingService:       istioConf.SkywalkingService,
			SkywalkingPort:          istioConf.SkywalkingPort,
			SkywalkingAccessToken:   istioConf.SkywalkingAccessToken,
			EnableZipkin:            istioConf.EnableZipkin,
			ZipkinService:           istioConf.ZipkinService,
			ZipkinPort:              istioConf.ZipkinPort,
			AccessLogServiceAddress: istioConf.AccessLogServiceAddress,
			MetricsServiceAddress:   istioConf.MetricsServiceAddress,
		}, local)
	case define.NvidiaRuntime:
		return release.RenderNvidiaTemplate(&release.NvidiaTemplate{
			MirrorUrl: current.Mirror,
		}, local)
	case define.KataRuntime:
		return release.RenderKataTemplate(&release.KataTemplate{
			MirrorUrl: current.Mirror,
		}, local)
	case define.KruisePlugin:
		return release.RenderKruiseTemplate(&release.KruiseTemplate{
			MirrorUrl: current.Mirror,
		}, local)
	case define.HealthzReader:
		return release.RenderHealthzTemplate(current.Version.Full), nil
	case define.LocalHaproxy:
		return release.RenderHaproxyTemplate(&release.BalanceHaproxyTemplate{
			MasterHosts: current.MasterAPIs(),
			ImageUrl:    current.GetHaproxyResource() + ":" + current.Version.Full,
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
			ImageUrl:    current.GetUpdaterResource() + ":" + current.Version.Full,
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

func calicoMode(conf *cluster.CalicoConf, dual bool) (string, string, string, string) {
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
		if dual {
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

func ciliumMode(conf *cluster.CiliumConf, dual bool) (bool, string, bool, string, string, bool) {
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
		ipv4NRCIDR = conf.NativeRoutingCIDR
		if dual {
			ipv6NRCIDR = conf.NativeRoutingCIDRV6
		}
		tunnelMode = define.CiliumTunnelDisabled
		autoDirectNode = true
	}
	return true, ipv4NRCIDR, dual, ipv6NRCIDR, tunnelMode, autoDirectNode
}

func ciliumLBMode(conf *cluster.CiliumConf) (string, string, string) {
	lbMode := define.CiliumLBModeSNAT
	if conf.EnableNR {
		lbMode = define.CiliumLBModeHybrid
	}
	if conf.EnableDSR {
		lbMode = define.CiliumLBModeDSR
	}
	lbAcceleration := ""
	if conf.LBNativeAcceleration {
		lbAcceleration = define.CiliumLBAccelerationNative
	}
	lbAlgorithm := ""
	if conf.LBMaglevAlgorithm {
		lbAlgorithm = define.CiliumLBAlgorithmMaglev
	}
	return lbMode, lbAcceleration, lbAlgorithm
}
