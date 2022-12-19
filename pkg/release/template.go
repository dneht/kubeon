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

package release

import (
	"bytes"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/release/configset"
	"github.com/pkg/errors"
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
	MirrorUrl    string
	ClusterDnsIP string
}

type CalicoTemplate struct {
	MirrorUrl              string
	IsSetInterface         bool
	DefaultInterface       string
	BackendMode            string
	EnableBPF              bool
	EnableWireGuard        bool
	CalicoMTU              uint32
	LBMode                 string
	EnableVXLAN            bool
	IPIPMode               string
	VXLANMode              string
	VXLANv6Mode            string
	BPFHostConntrackBypass bool
	ClusterEnableDual      bool
	ClusterLbDomain        string
	ClusterLbPort          uint32
	ClusterPortRange       string
	ClusterNodeMaskSize    uint32
	ClusterNodeMaskSizeV6  uint32
	ClusterPodCIDR         string
	ClusterPodCIDRV6       string
}

type CiliumTemplate struct {
	MirrorUrl                       string
	IsSetInterface                  bool
	DefaultInterface                string
	EnableBGP                       bool
	EnableBM                        bool
	EnableBBR                       bool
	EnableWireGuard                 bool
	EnableIPv4Masquerade            bool
	EnableIPv6Masquerade            bool
	NativeRoutingCIDR               string
	NativeRoutingCIDRV6             string
	EnableIPv6BigTCP                bool
	CiliumMTU                       uint32
	TunnelMode                      string
	PolicyMode                      string
	LBMode                          string
	LBAcceleration                  string
	LBAlgorithm                     string
	LBHostNamespaceOnly             bool
	AutoDirectNodeRoutes            bool
	EnableLocalRedirect             bool
	AutoProtectPortRange            bool
	BPFMapDynamicSizeRatio          string
	BPFLBMapMax                     uint32
	BPFPolicyMapMax                 uint32
	BPFLBExternalClusterIP          bool
	BPFLBBypassFIBLookup            bool
	InstallIptablesRules            bool
	InstallNoConntrackIptablesRules bool
	MonitorAggregation              string
	MonitorInterval                 string
	MonitorFlags                    string
	ClusterEnableDual               bool
	ClusterLbDomain                 string
	ClusterLbPort                   uint32
	ClusterPortRange                string
	ClusterNodeMaskSize             uint32
	ClusterNodeMaskSizeV6           uint32
	ClusterPodCIDR                  string
	ClusterPodCIDRV6                string
}

type NvidiaTemplate struct {
	MirrorUrl string
}

type KataTemplate struct {
	MirrorUrl string
}

type ContourTemplate struct {
	MirrorUrl             string
	DisableMergeSlashes   bool
	DisablePermitInsecure bool
}

type IstioTemplate struct {
	MirrorUrl               string
	EnableNetworkPlugin     bool
	ProxyClusterDomain      string
	EnableAutoInject        bool
	ServiceEntryExportTo    []string
	VirtualServiceExportTo  []string
	DestinationRuleExportTo []string
	EnableAutoMTLS          bool
	EnableHttp2AutoUpgrade  bool
	EnablePrometheusMerge   bool
	EnableIngressGateway    bool
	IngressGatewayType      string
	EnableEgressGateway     bool
	EgressGatewayType       string
	EnableSkywalking        bool
	EnableSkywalkingAll     bool
	SkywalkingService       string
	SkywalkingPort          uint32
	SkywalkingAccessToken   string
	EnableZipkin            bool
	ZipkinService           string
	ZipkinPort              uint32
	AccessLogServiceAddress string
	MetricsServiceAddress   string
}

type KruiseTemplate struct {
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
	if local {
		return parseFile("/cilium.yaml.tpl", input)
	} else {
		return parseFile("/cilium.m.yaml.tpl", input)
	}
}

func RenderContourTemplate(input *ContourTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/contour.yaml.tpl", input)
	} else {
		return parseFile("/contour.m.yaml.tpl", input)
	}
}

func RenderIstioTemplate(input *IstioTemplate, local bool) ([]byte, error) {
	istioArgs := BuildIstioctlArgs(input, false, local)
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
