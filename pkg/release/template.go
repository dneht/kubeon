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
	APIVersion          string
	ImageRepository     string
	ClusterName         string
	ClusterVersion      string
	ClusterPortRange    string
	ClusterFeatureGates string
	ClusterApiIP        string
	ClusterLbIP         string
	ClusterLbPort       int32
	ClusterLbDomain     string
	ClusterDnsDomain    string
	ClusterSvcCIDR      string
	ClusterPodCIDR      string
	IsExternalLB        bool
	MasterCertSANs      []string
	InputCertSANs       []string
	KubeProxyMode       string
	KubeIPVSScheduler   string
}

type KubeadmInitTemplate struct {
	APIVersion       string
	Token            string
	NodeName         string
	ImagePullPolicy  string
	AdvertiseAddress string
	BindPort         int32
	CertificateKey   string
}

type KubeadmJoinTemplate struct {
	APIVersion       string
	Token            string
	ClusterLbDomain  string
	ClusterLbPort    int32
	CaCertHash       string
	IsControlPlane   bool
	NodeName         string
	ImagePullPolicy  string
	AdvertiseAddress string
	BindPort         int32
	CertificateKey   string
}

type CorednsTemplate struct {
	ClusterDnsIP string
}

type CalicoTemplate struct {
	IsSetInterface   bool
	DefaultInterface string
	EtcdKeyBase64    string
	EtcdCertBase64   string
	EtcdCABase64     string
	EtcdEndpoints    string
	CalicoMTU        string
	IPIPMode         string
	VXLanMode        string
}

type NvidiaTemplate struct {
}

type KataTemplate struct {
}

type ContourTemplate struct {
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

func RenderContourTemplate(input *ContourTemplate, local bool) ([]byte, error) {
	if local {
		return parseFile("/contour.yaml.tpl", input)
	} else {
		return parseFile("/contour.m.yaml.tpl", input)
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
