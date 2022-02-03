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
	}
	return nil
}

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
	switch current.IngressMode {
	case define.ContourIngress:
		{
			err = InstallInner(define.ContourIngress)
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

func getCorednsTemplate() *release.CorednsTemplate {
	current := cluster.Current()
	return &release.CorednsTemplate{
		ClusterDnsIP: current.DnsIP,
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
		return release.RenderCorednsTemplate(getCorednsTemplate(), local)
	case define.CalicoNetwork:
		etcdConfig := current.CreateConfig
		if nil == etcdConfig {
			return nil, errors.New("get etcd config error")
		}
		ipipMode := "Always"
		vxlanMode := "Never"
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
		}, local)
	case define.NvidiaRuntime:
		return release.RenderNvidiaTemplate(&release.NvidiaTemplate{}, local)
	case define.KataRuntime:
		return release.RenderKataTemplate(&release.KataTemplate{}, local)
	case define.ContourIngress:
		return release.RenderContourTemplate(&release.ContourTemplate{}, local)
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
