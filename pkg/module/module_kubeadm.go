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
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"time"
)

var currentToken string
var currentSecret string

func SetupBootKubeadm(node *cluster.Node) (err error) {
	current := cluster.Current()
	nodes := cluster.NodeList{node}
	err = generateNodeKubeadm(current, nodes, false)
	if nil != err {
		return err
	}
	err = sendNeedKubeadm(nodes)
	if nil != err {
		return err
	}
	return action.KubeadmInitStart(node, true, false, define.KubeadmIgnorePreflightErrors)
}

func SetupJoinsKubeadm(nodes cluster.NodeList) (err error) {
	current := cluster.Current()
	err = generateNodeKubeadm(current, nodes, false)
	if nil != err {
		return err
	}
	err = sendNeedKubeadm(nodes)
	if nil != err {
		return err
	}
	return action.KubeadmJoinNode(nodes, false, define.KubeadmIgnorePreflightErrors, 2*time.Minute)
}

func SetupUpgradeKubeadm() (err error) {
	current := cluster.Current()
	currentNodes := cluster.CurrentNodes()
	err = generateNodeKubeadm(current, currentNodes, true)
	if nil != err {
		return err
	}
	return sendNeedKubeadm(currentNodes)
}

func SetupAddKubeadm(nodes cluster.NodeList) (err error) {
	current := cluster.Current()
	err = generateNodeKubeadm(current, cluster.CurrentNodes(), true)
	if nil != err {
		return err
	}
	err = sendNeedKubeadm(nodes)
	if nil != err {
		return err
	}
	return action.KubeadmJoinNode(nodes, false, define.KubeadmIgnorePreflightErrors, 2*time.Minute)
}

func sendNeedKubeadm(nodes cluster.NodeList) (err error) {
	for _, node := range nodes {
		nodeConf := node.GetResource().ClusterConf
		err = node.CopyTo(node.LocalConfigPath(), nodeConf.KubeadmInitPath)
		if nil != err {
			return err
		}
	}
	return nil
}

func generateNodeKubeadm(current *cluster.Cluster, nodes cluster.NodeList, needCreate bool) (err error) {
	localConf := cluster.CurrentResource().ClusterConf
	onutil.MkDir(localConf.KubeadmInitDir)
	if "" == currentToken {
		currentToken, err = cluster.GenerateToken()
		if nil != err {
			return err
		}
	}
	if "" == currentSecret {
		currentSecret = onutil.GetSecretSHA265()
	}
	for _, node := range nodes {
		err = writeKubeadmConfig(current, node, currentToken, currentSecret)
		if nil != err {
			return err
		}
	}
	if needCreate {
		err = action.KubeadmUploadCert(currentSecret)
		if nil != err {
			return err
		}
		return cluster.CreateToken()
	} else {
		return nil
	}
}

func writeKubeadmConfig(current *cluster.Cluster, node *cluster.Node, token, secretKey string) (err error) {
	if nil == current {
		return errors.New("cluster not init")
	}
	cps := current.ControlPlanes
	masterCertSANs := make([]string, len(cps)*2)
	for i, cp := range cps {
		n := i * 2
		masterCertSANs[n] = cp.IPv4
		masterCertSANs[n+1] = cp.Hostname
	}
	if nil == current.CertSANs {
		current.CertSANs = []string{}
	}
	kubeletTemplate := getKubeletTemplate()
	kubeadmTemplate := &release.KubeadmTemplate{
		ClusterName:       current.Name,
		ClusterVersion:    current.Version.Full,
		ClusterPortRange:  current.PortRange,
		ClusterApiIP:      current.ApiIP,
		ClusterLbIP:       current.LbIP,
		ClusterLbPort:     current.LbPort,
		ClusterLbDomain:   current.LbDomain,
		ClusterDnsDomain:  current.DnsDomain,
		ClusterSvcCIDR:    current.SvcCIDR,
		ClusterPodCIDR:    current.PodCIDR,
		IsExternalLB:      current.IsExternalLb,
		MasterCertSANs:    masterCertSANs,
		InputCertSANs:     current.CertSANs,
		KubeProxyMode:     current.ProxyMode,
		KubeIPVSScheduler: current.IPVSScheduler,
	}
	if node.IsBootstrap() {
		kubeadmTemplate.ClusterLbPort = define.DefaultClusterAPIPort
		err = release.WriteKubeadmInitTemplate(kubeletTemplate, kubeadmTemplate,
			&release.KubeadmInitTemplate{
				Token:            token,
				NodeName:         node.Hostname,
				AdvertiseAddress: node.IPv4,
				BindPort:         define.DefaultClusterAPIPort,
				CertificateKey:   secretKey,
			}, current.Version.Full, node.LocalConfigPath())
	} else {
		lbPort := current.LbPort
		if node.IsControlPlane() {
			lbPort = define.DefaultClusterAPIPort
		}
		kubeadmTemplate.ClusterLbPort = lbPort
		err = release.WriteKubeadmJoinTemplate(kubeletTemplate, kubeadmTemplate,
			&release.KubeadmJoinTemplate{
				Token:            token,
				ClusterLbDomain:  current.LbDomain,
				ClusterLbPort:    lbPort,
				CaCertHash:       current.GetCertHash(),
				IsControlPlane:   node.IsControlPlane(),
				NodeName:         node.Hostname,
				AdvertiseAddress: node.IPv4,
				BindPort:         define.DefaultClusterAPIPort,
				CertificateKey:   secretKey,
			}, current.Version.Full, node.LocalConfigPath())
	}
	return err
}
