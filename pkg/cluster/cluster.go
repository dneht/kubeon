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

package cluster

import (
	"encoding/base64"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"strconv"
)

var current *Cluster

type Cluster struct {
	Name                 string                   `json:"name"`
	Version              *define.StdVersion       `json:"version"`
	Mirror               string                   `json:"mirror,omitempty"`
	IsBinary             bool                     `json:"binary"`
	IsOffline            bool                     `json:"offline"`
	UsePatch             bool                     `json:"usePatch"`
	ApiIP                string                   `json:"apiIP"`
	DnsIP                string                   `json:"dnsIP"`
	LbIP                 string                   `json:"lbIP"`
	LbPort               int32                    `json:"lbPort"`
	LbDomain             string                   `json:"lbDomain"`
	LbMode               string                   `json:"lbMode"`
	DnsDomain            string                   `json:"dnsDomain"`
	MaxPods              uint32                   `json:"maxPods"`
	PortRange            string                   `json:"portRange"`
	SvcCIDR              string                   `json:"svcCIDR"`
	PodCIDR              string                   `json:"podCIDR"`
	NodeInterface        []string                 `json:"nodeInterface"`
	ProxyMode            string                   `json:"proxyMode"`
	IPVSScheduler        string                   `json:"ipvsScheduler"`
	RuntimeMode          string                   `json:"runtimeMode"`
	NetworkMode          string                   `json:"networkMode"`
	CalicoMode           string                   `json:"calicoMode"`
	CalicoMTU            string                   `json:"calicoMTU"`
	IngressMode          string                   `json:"ingressMode"`
	UseNvidia            bool                     `json:"useNvidia"`
	HasNvidia            bool                     `json:"hasNvidia"`
	UseKata              bool                     `json:"useKata"`
	ControlPlanes        NodeList                 `json:"controlPlanes"`
	Workers              NodeList                 `json:"workers"`
	AllNodes             NodeList                 `json:"-"`
	IsExternalLb         bool                     `json:"isExternalLb"`
	IsExternalEtcd       bool                     `json:"isExternalEtcd"`
	CertSANs             []string                 `json:"certSANs,omitempty"`
	CreateConfig         *CreateConfig            `json:"createConfig,omitempty"`
	LocalResource        *release.ClusterResource `json:"localResource,omitempty"`
	ExistResourceVersion *map[string]string       `json:"-"`
	AdminConfigPath      string                   `json:"adminConfigPath,omitempty"`
	Status               RunStatus                `json:"status"`
}

type RunStatus string

const (
	StatusCreating   RunStatus = "creating"
	StatusUpgrading  RunStatus = "upgrading"
	StatusAddWaiting RunStatus = "addWaiting"
	StatusDelWaiting RunStatus = "delWaiting"
	StatusRunning    RunStatus = "running"
)

func Current() *Cluster {
	return current
}

func CurrentNodes() NodeList {
	return current.AllNodes
}

func CurrentResource() *release.ClusterResource {
	return current.LocalResource
}

func GenerateToken() (string, error) {
	result, err := execute.NewLocalCmd("kubeadm", "token", "generate").RunAndResult()
	if nil != err {
		return "", err
	}
	return result, nil
}

func CreateToken() error {
	return execute.NewLocalCmd("kubeadm", "token", "create",
		"--config="+BootstrapNode().LocalConfigPath()).Run()
}

func (c *Cluster) AllIPs() []string {
	all := make([]string, len(c.AllNodes))
	for i, node := range c.AllNodes {
		all[i] = node.IPv4
	}
	return all
}

func (c *Cluster) MasterIPs() []string {
	all := make([]string, len(current.ControlPlanes))
	for i, node := range current.ControlPlanes {
		all[i] = node.IPv4
	}
	return all
}

func (c *Cluster) MasterAPIs() []string {
	all := make([]string, len(current.ControlPlanes))
	port := strconv.FormatInt(int64(define.DefaultClusterAPIPort), 10)
	for i, node := range current.ControlPlanes {
		all[i] = node.IPv4 + ":" + port
	}
	return all
}

func (c *Cluster) Healthz() string {
	return "https://" + c.LbDomain + ":" + strconv.FormatInt(int64(c.LbPort), 10) + "/healthz"
}

func (c *Cluster) InCluster(ipv4 string) bool {
	if nil == c {
		return false
	}
	for _, node := range c.ControlPlanes {
		if node.IPv4 == ipv4 {
			return true
		}
	}
	for _, node := range c.Workers {
		if node.IPv4 == ipv4 {
			return true
		}
	}
	return false
}

func (c *Cluster) IsMultiMaster() bool {
	return len(c.ControlPlanes) > 1
}

func (c *Cluster) HasPureWorker() bool {
	return len(c.Workers) > 0
}

func (c *Cluster) GetCertHash() (string, error) {
	if nil == c {
		return "", errors.New("cluster not found")
	}
	if nil == c.CreateConfig {
		loadCreateConfig()
	}
	if nil == c.CreateConfig {
		return "", errors.New("cluster create config get error")
	}
	caData, err := base64.StdEncoding.DecodeString(c.CreateConfig.CACertBase64)
	if nil != err {
		return "", err
	}
	return "sha256:" + onutil.CertSHA256(caData), nil
}

func (c *Cluster) IsRealLocal() bool {
	return c.IsOffline || c.Version.LessThen(define.K8S_1_22_0)
}

func (c *Cluster) GetInitImageRepo() string {
	if c.IsRealLocal() {
		return define.DefaultImageRepo
	} else {
		return define.MirrorImageRepo + "/kubeon"
	}
}

func (c *Cluster) GetInitImagePullPolicy() string {
	if c.IsRealLocal() {
		return "IfNotPresent"
	} else {
		return "Always"
	}
}

func (c *Cluster) GetHaproxyResource() string {
	if c.IsRealLocal() {
		return define.HaproxyResource
	} else {
		return define.MirrorImageRepo + "/" + define.HaproxyResource
	}
}

func (c *Cluster) GetUpdaterResource() string {
	if c.IsRealLocal() {
		return define.UpdaterResource
	} else {
		return define.MirrorImageRepo + "/" + define.UpdaterResource
	}
}

func (c *Cluster) GetExistVer(mod string) string {
	return (*c.ExistResourceVersion)[mod]
}

func (c *Cluster) ModuleVersionChange(mod, iVer string) bool {
	if len(mod) == 0 || nil == c.ExistResourceVersion {
		return true
	}
	eVer := (*c.ExistResourceVersion)[mod]
	return eVer != iVer
}

func (c *Cluster) GetKubeletAPIVersion() string {
	return define.KubeletConfigApiB1
}

func (c *Cluster) GetKubeadmAPIVersion() string {
	if c.Version.LessThen(define.K8S_1_22_0) {
		return define.KubeadmConfigApiB2
	} else {
		return define.KubeadmConfigApiB3
	}
}

func (c *Cluster) GetKubeadmFeatureGates() string {
	if c.Version.LessThen(define.K8S_1_23_0) {
		return "feature-gates: TTLAfterFinished=true"
	} else {
		return ""
	}
}
