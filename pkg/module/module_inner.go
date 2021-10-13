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
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"strings"
)

func InstallInner(moduleName string) (err error) {
	bytes, err := ShowInner(moduleName)
	if nil != err {
		return err
	}
	if nil != bytes {
		log.Debugf("install %s on cluster", moduleName)
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
		ClusterDnsIP:     current.DnsIP,
		ClusterDnsDomain: current.DnsDomain,
		ClusterMaxPods:   current.MaxPods,
	}
}

func getCorednsTemplate() *release.CorednsTemplate {
	current := cluster.Current()
	return &release.CorednsTemplate{
		ClusterDnsIP: current.DnsIP,
	}
}

func ShowInner(moduleName string) ([]byte, error) {
	current := cluster.Current()
	switch moduleName {
	case define.KubeletModule:
		log.Debugf("get module %s config", define.KubeletModule)
		return release.RenderKubeletTemplate(getKubeletTemplate(), current.Version.Full)
	case define.CorednsPart:
		log.Debugf("get module %s config", define.CorednsPart)
		return release.RenderCorednsTemplate(getCorednsTemplate())
	case define.CalicoNetwork:
		log.Debugf("get module %s config", define.CalicoNetwork)
		ipipMode := "Always"
		vxlanMode := "Never"
		etcdConfig := current.CreateConfig
		if nil == etcdConfig {
			return nil, errors.New("get etcd config error")
		}
		if define.CalicoVXLan == current.CalicoMode {
			ipipMode = "Never"
			vxlanMode = "Always"
		}
		nodeInterface := current.NodeInterface
		isSetInterface := len(nodeInterface) != 0
		defaultInterface := ""
		if isSetInterface {
			defaultInterface = nodeInterface[0]
		}
		return release.RenderCalicoTemplate(&release.CalicoTemplate{
			IsSetInterface:   isSetInterface,
			DefaultInterface: defaultInterface,
			EtcdKeyBase64:    etcdConfig.EtcdKeyBase64,
			EtcdCertBase64:   etcdConfig.EtcdCertBase64,
			EtcdCABase64:     etcdConfig.EtcdCABase64,
			EtcdEndpoints:    etcdConfig.EtcdEndpoints,
			CalicoMTU:        current.CalicoMTU,
			IPIPMode:         ipipMode,
			VXLanMode:        vxlanMode,
		})
	case define.ContourIngress:
		log.Debugf("get module %s config", define.ContourIngress)
		return release.RenderContourTemplate(&release.ContourTemplate{
		})
	case define.HealthzReader:
		log.Debugf("get module %s config", define.HealthzReader)
		return release.RenderHealthzTemplate(current.Version.Full), nil
	case define.LocalHaproxy:
		log.Debugf("get module %s config", define.LocalHaproxy)
		return release.RenderHaproxyTemplate(&release.BalanceHaproxyTemplate{
			MasterHosts: current.MasterAPIs(),
			ImageUrl:    define.HaproxyResource + ":" + current.Version.Full,
		})
	case define.ApiserverStartup:
		log.Debugf("get module %s config", define.ApiserverStartup)
		return release.RenderStartupService(&release.ApiserverScriptTemplate{
			TargetDomain: current.LbDomain,
			VirtualAddr:  current.LbIP,
			RealAddrs:    strings.Join(current.MasterIPs(), ","),
		})
	case define.ApiserverUpdater:
		log.Debugf("get module %s config", define.ApiserverUpdater)
		return release.RenderUpdaterTemplate(&release.ApiserverUpdaterTemplate{
			ClusterLbIP: current.LbIP,
			MasterIPs:   current.MasterIPs(),
			ImageUrl:    define.UpdaterResource + ":" + current.Version.Full,
		})
	default:
		log.Warnf("not support inner module[%s]", moduleName)
		return nil, nil
	}
}
