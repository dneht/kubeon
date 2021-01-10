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

package cluster

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"net"
	"strings"
)

func InitNewCluster(cluster *Cluster, lb string, base define.DefaultList, master define.MasterList, worker define.WorkerList) error {
	current = cluster
	if current.Status != StatusCreating && runConfig.Exist {
		configRaw, err := runConfig.ReadConfig()
		if nil != err {
			return err
		}
		fmt.Println(configRaw)
		return errors.Errorf("input cluster[%s] exist", runConfig.Name)
	}

	err, apiIP, dnsIP, isExtendLB, lbIP := parseCIDR(current.SvcCIDR, current.PodCIDR, lb, current.LbMode)
	if nil != err {
		return err
	}
	current.Name = runConfig.Name
	masterNodes, workerNodes, isReady := newNodeList(base, master, worker)
	if !isReady {
		return errors.New("init new cluster node error")
	}
	current.ApiIP = apiIP
	current.DnsIP = dnsIP
	current.LbIP = lbIP
	current.IsExternalLb = isExtendLB
	if !isExtendLB {
		//if !current.IsMultiMaster() {
		//	current.LbIP = masterNodes[0].IPv4
		//}
		current.LbPort = define.DefaultClusterAPIPort
	}
	current.ControlPlanes = masterNodes
	current.Workers = workerNodes
	current.LocalResource = release.InitResource(current.Version.Full, current.RuntimeMode, current.IsBinary, current.IsOffline)
	current.AdminConfigPath = define.AppConfDir + "/cluster/" + current.Name + ".yaml"

	initCurrent(current)
	err = InitHost()
	if nil != err {
		log.Errorf("init node host error: %v", err)
		return err
	}
	return runConfig.WriteConfig()
}

func InitExistCluster() (*RunConfig, error) {
	if nil == runConfig || !runConfig.Exist {
		return nil, errors.Errorf("config[%s] not exist, please create", runConfig.Name)
	} else {
		cluster, err := runConfig.ParseConfig()
		if nil != err {
			return nil, err
		}
		initCurrent(cluster)
		for _, node := range currentNodes {
			node.SetConnect()
		}
	}
	return runConfig, nil
}

func InitUpgradeCluster(version *define.StdVersion) error {
	if current.Status != StatusUpgrading && version.LessThen(current.Version) {
		return errors.Errorf("upgrade version [%s] is less than now version [%s]", version.Full, current.Version.Full)
	}

	current.Version = version
	current.LocalResource = release.InitResource(version.Full, current.RuntimeMode, current.IsBinary, current.IsOffline)
	current.Status = StatusUpgrading

	err := InitHost()
	if nil != err {
		log.Errorf("init node host error: %v", err)
		return err
	}
	return runConfig.WriteConfig()
}

func AfterBuildCluster() (*CreateConfig, error) {
	var err error
	adminConfigBase64 := bootBase64(define.KubernetesEtcPath + "/admin.conf")
	current.CreateConfig = &CreateConfig{
		CACertBase64:      bootBase64(define.KubernetesPkiPath + "/ca.crt"),
		EtcdKeyBase64:     bootBase64(define.KubernetesPkiPath + "/apiserver-etcd-client.key"),
		EtcdCertBase64:    bootBase64(define.KubernetesPkiPath + "/apiserver-etcd-client.crt"),
		EtcdCABase64:      bootBase64(define.KubernetesPkiEtcdPath + "/ca.crt"),
		EtcdEndpoints:     etcdEndpoints(),
		AdminConfigBase64: adminConfigBase64,
	}
	err = writeKubeConfig(adminConfigBase64)
	if err != nil {
		return nil, err
	}
	err = runConfig.ChangeConfig()
	if nil != err {
		log.Error("save cluster config failed: " + err.Error())
	}
	return current.CreateConfig, err
}

func InitAddNodes(base define.DefaultList, master define.MasterList, worker define.WorkerList) (NodeList, error) {
	masterNodes, workerNodes, isReady := newNodeList(base, master, worker)
	if !isReady {
		return nil, errors.New("init add cluster node error")
	}
	newNodes := MergeNodeList(masterNodes, workerNodes)
	err := checkExist(newNodes)
	if nil != err {
		return nil, err
	}

	current.ControlPlanes = MergeNodeList(current.ControlPlanes, masterNodes)
	current.Workers = MergeNodeList(current.Workers, workerNodes)
	if nil != current.CreateConfig {
		current.CreateConfig.EtcdEndpoints = etcdEndpoints()
	}

	initCurrent(current)
	err = InitHost()
	if nil != err {
		log.Errorf("init node host error: %v", err)
		return nil, err
	}
	return newNodes, runConfig.WriteConfig()
}

func InitDelNodes(selector string) (NodeList, error) {
	delNodes, err := SelectNodes(selector)
	if nil != err {
		return nil, err
	}
	if len(delNodes) == 0 {
		return nil, errors.Errorf("selector[%s] can not get some node", selector)
	}
	hash := make(map[string]*Node, len(currentNodes))
	for _, node := range currentNodes {
		hash[node.IPv4] = node
	}
	for _, node := range delNodes {
		delete(hash, node.IPv4)
	}

	controlPlanes := make(NodeList, 0, len(delNodes))
	workers := make(NodeList, 0, len(delNodes))
	for _, node := range hash {
		if node.IsControlPlane() {
			controlPlanes = append(controlPlanes, node)
		} else if node.IsWorker() {
			workers = append(workers, node)
		}
	}
	cpSize := len(controlPlanes)
	if cpSize == 0 {
		return nil, errors.Errorf("selector[%s] no master exists after delete", selector)
	}
	current.ControlPlanes = SortNodeList(controlPlanes)
	current.Workers = SortNodeList(workers)
	if nil != current.CreateConfig {
		current.CreateConfig.EtcdEndpoints = etcdEndpoints()
	}

	initCurrent(current)
	return delNodes, runConfig.WriteConfig()
}

func initCurrent(cluster *Cluster) {
	currentNodes = MergeNodeList(cluster.ControlPlanes, cluster.Workers)
}

func parseCIDR(svcCIDR, podCIDR, lbIP, lbMode string) (error, string, string, bool, string) {
	start, _, err := net.ParseCIDR(svcCIDR)
	if err != nil {
		log.Errorf("cluster service CIDR error, please check [%s]", svcCIDR)
		return err, "", "", false, ""
	}
	_, _, err = net.ParseCIDR(podCIDR)
	if err != nil {
		log.Errorf("cluster pod CIDR error, please check [%s]", podCIDR)
		return err, "", "", false, ""
	}
	start = start.To4()
	start[3] = 1
	api := start.String()
	start[3] += 9
	dns := start.String()
	extend := true
	lb := lbIP
	if "" == lb {
		extend = false
		if lbMode == define.LocalHaproxy {
			lb = "127.0.0.1"
		} else {
			start[3] += 5
			lb = start.String()
		}
	}
	return nil, api, dns, extend, lb
}

func etcdEndpoints() string {
	endpoints := make([]string, len(current.ControlPlanes))
	for idx, cp := range current.ControlPlanes {
		endpoints[idx] = cp.IPv4 + ":2379"
	}
	return strings.Join(endpoints, ",")
}
