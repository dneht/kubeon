/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package release

import (
	"bytes"
	"fmt"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/release/configset"
	"github.com/pkg/errors"
	"strings"
	"text/template"
)

type KubeletTemplate struct {
	APIVersion       string
	ClusterDnsIP     string
	ClusterDnsDomain string
	ClusterMaxPods   uint32
}

type KubeadmTemplate struct {
	APIVersion             string
	ImageRepository        string
	ClusterName            string
	ClusterVersion         string
	ClusterEnableDual      bool
	ClusterPortRange       string
	ClusterNodeMaskSize    uint32
	ClusterNodeMaskSizeV6  uint32
	ClusterFeatureGates    string
	ClusterSigningDuration string
	ClusterApiIP           string
	ClusterLbIP            string
	ClusterLbPort          uint32
	ClusterLbDomain        string
	ClusterDnsDomain       string
	ClusterSvcCIDR         string
	ClusterPodCIDR         string
	IsExternalLB           bool
	MasterCertSANs         []string
	InputCertSANs          []string
	ProxyMode              string
	IPVSScheduler          string
	StrictARP              bool
}

type KubeadmInitTemplate struct {
	APIVersion       string
	Token            string
	NodeName         string
	ImagePullPolicy  string
	AdvertiseAddress string
	BindPort         uint32
	CertificateKey   string
}

type KubeadmJoinTemplate struct {
	APIVersion       string
	Token            string
	ClusterLbDomain  string
	ClusterLbPort    uint32
	CaCertHash       string
	IsControlPlane   bool
	NodeName         string
	ImagePullPolicy  string
	AdvertiseAddress string
	BindPort         uint32
	CertificateKey   string
}

type CorednsTemplate struct {
	CPVersion    string
	MirrorUrl    string
	ClusterDnsIP string
}

type CalicoTemplate struct {
	CPVersion             string
	MirrorUrl             string
	IsSetInterface        bool
	DefaultInterface      string
	BackendMode           string
	EnableBPF             bool
	EnableWireGuard       bool
	CalicoMTU             uint32
	LBMode                string
	EnableVXLAN           bool
	IPIPMode              string
	VXLANMode             string
	VXLANv6Mode           string
	EnablePassConntrack   bool
	ClusterEnableDual     bool
	ClusterLbDomain       string
	ClusterLbPort         uint32
	ClusterPortRange      string
	ClusterNodeMaskSize   uint32
	ClusterNodeMaskSizeV6 uint32
	ClusterPodCIDR        string
	ClusterPodCIDRV6      string
}

type CiliumTemplate struct {
	CPVersion               string
	MirrorUrl               string
	IsSetInterface          bool
	DefaultInterface        string
	EnableBGP               bool
	EnableBM                bool
	EnableBBR               bool
	EnableWireGuard         bool
	EnableIPv4Masquerade    bool
	EnableIPv6Masquerade    bool
	NativeRoutingCIDR       string
	NativeRoutingCIDRV6     string
	CiliumMTU               uint32
	PolicyMode              string
	TunnelMode              string
	LBMode                  string
	EnableEndpointRoutes    bool
	EnableLocalRedirect     bool
	EnableHostnsOnly        bool
	AutoDirectNodeRoutes    bool
	EnableEndpointSlice     bool
	EnableExternalClusterIP bool
	AutoProtectPortRange    bool
	HubbleVersion           string
	EnableHubbleTLS         bool
	CustomConfigs           []string
	ClusterEnableDual       bool
	ClusterLbDomain         string
	ClusterLbPort           uint32
	ClusterPortRange        string
	ClusterNodeMaskSize     uint32
	ClusterNodeMaskSizeV6   uint32
	ClusterPodCIDR          string
	ClusterPodCIDRV6        string
}

type NvidiaTemplate struct {
	CPVersion string
	MirrorUrl string
}

type KataTemplate struct {
	CPVersion string
	MirrorUrl string
}

type ContourTemplate struct {
	CPVersion             string
	MirrorUrl             string
	DisableMergeSlashes   bool
	DisablePermitInsecure bool
}

type IstioTemplate struct {
	CPVersion               string
	MirrorUrl               string
	EnableAutoInject        bool
	ServiceEntryExportTo    []string
	VirtualServiceExportTo  []string
	DestinationRuleExportTo []string
	IngressGatewayType      string
	EnableEgressGateway     bool
	EgressGatewayType       string
	CustomConfigs           []string
	ProxyClusterDomain      string
}

type KruiseTemplate struct {
	CPVersion    string
	MirrorUrl    string
	FeatureGates string
}

type BalanceHaproxyTemplate struct {
	MasterHosts []string
	ImageUrl    string
}

type ApiserverScriptTemplate struct {
	TargetDomain string
	VirtualAddr  string
	RealAddrs    string
}

type ApiserverUpdaterTemplate struct {
	ClusterLbIP string
	MasterIPs   []string
	ImageUrl    string
}

func RenderCorednsTemplate(input *CorednsTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/coredns.yaml.tpl", input)
	} else {
		return parseFile("/coredns.m.yaml.tpl", input)
	}
}

func RenderCalicoTemplate(input *CalicoTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/calico.yaml.tpl", input)
	} else {
		return parseFile("/calico.m.yaml.tpl", input)
	}
}

func RenderCiliumTemplate(input *CiliumTemplate, local bool) ([]byte, error) {
	ciliumArgs := buildCiliumInstallArgs(input, false, local)
	hubbleArgs := buildHubbleInstallArgs(input, false, local)
	return []byte(fmt.Sprintf("%s %s\n%s %s",
		define.CiliumCommand, strings.Join(ciliumArgs, " "),
		define.CiliumCommand, strings.Join(hubbleArgs, " "))), nil
}

func RenderContourTemplate(input *ContourTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/contour.yaml.tpl", input)
	} else {
		return parseFile("/contour.m.yaml.tpl", input)
	}
}

func RenderIstioTemplate(input *IstioTemplate, local bool) ([]byte, error) {
	istioArgs := buildIstioInstallArgs(input, false, local)
	istioCmd := execute.NewLocalCmd(define.IstioCommand, istioArgs...)
	return istioCmd.RunAndBytes()
}

func RenderNvidiaTemplate(input *NvidiaTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/nvidia.yaml.tpl", input)
	} else {
		return parseFile("/nvidia.m.yaml.tpl", input)
	}
}

func RenderKataTemplate(input *KataTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/kata.yaml.tpl", input)
	} else {
		return parseFile("/kata.m.yaml.tpl", input)
	}
}

func RenderKruiseTemplate(input *KruiseTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/kruise.yaml.tpl", input)
	} else {
		return parseFile("/kruise.m.yaml.tpl", input)
	}
}

func parseFile(path string, input interface{}) ([]byte, error) {
	tplPath := define.AppTplDir + path
	if !onutil.PathExists(tplPath) {
		return nil, errors.Errorf("Can not found template[%s]", tplPath)
	}
	tmpl, err := template.ParseFiles(tplPath)
	if nil != err {
		return nil, err
	}
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, input)
	return buffer.Bytes(), nil
}

func RenderKubeletTemplate(full *KubeletTemplate, version string) ([]byte, error) {
	return parseInner(configset.GetKubeletByVersion(version), full)
}

func WriteKubeadmInitTemplate(svc *KubeletTemplate, base *KubeadmTemplate, init *KubeadmInitTemplate, version, initPath string) error {
	baseArr, err := parseInner(configset.GetKubeadmByVersion(version), base)
	if nil != err {
		return err
	}
	initArr, err := parseInner(configset.GetKubeadmInitByVersion(version), init)
	if nil != err {
		return err
	}
	svcArr, err := RenderKubeletTemplate(svc, version)
	if nil != err {
		return err
	}

	var buff bytes.Buffer
	buff.Write(initArr)
	buff.Write(baseArr)
	buff.Write(svcArr)
	return writeFile(initPath, buff.Bytes())
}

func WriteKubeadmJoinTemplate(svc *KubeletTemplate, base *KubeadmTemplate, join *KubeadmJoinTemplate, version, joinPath string) error {
	baseArr, err := parseInner(configset.GetKubeadmByVersion(version), base)
	if nil != err {
		return err
	}
	joinArr, err := parseInner(configset.GetKubeadmJoinByVersion(version), join)
	if nil != err {
		return err
	}
	svcArr, err := RenderKubeletTemplate(svc, version)
	if nil != err {
		return err
	}

	var buff bytes.Buffer
	buff.Write(joinArr)
	buff.Write(baseArr)
	buff.Write(svcArr)
	return writeFile(joinPath, buff.Bytes())
}

func RenderHealthzTemplate(version string) []byte {
	return []byte(configset.GetHealthzReaderByVersion(version))
}

func RenderHaproxyTemplate(input *BalanceHaproxyTemplate) ([]byte, error) {
	return parseInner(configset.GetHaproxyStaticTemplate(), input)
}

func RenderStartupService(input *ApiserverScriptTemplate) ([]byte, error) {
	return parseInner(configset.GetApiserverStartupScript(), input)
}

func WriteStartupService(arr []byte, svcPath, bashPath string) error {
	err := writeFile(svcPath, []byte(configset.GetApiserverStartupService()))
	if nil != err {
		return err
	}
	return writeFile(bashPath, arr)
}

func RenderUpdaterTemplate(input *ApiserverUpdaterTemplate) ([]byte, error) {
	return parseInner(configset.GetApiserverUpdaterTemplate(), input)
}

func parseInner(inner string, input interface{}) ([]byte, error) {
	tmpl, err := template.New("inner").Parse(inner)
	if nil != err {
		return nil, err
	}
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, input)
	return buffer.Bytes(), nil
}

func writeFile(path string, arr []byte) error {
	return onutil.WriteFile(path, arr)
}
