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
	"path/filepath"
	"time"
)

func InstallLoadBalance(nodes cluster.NodeList) (err error) {
	current := cluster.Current()
	if current.IsExternalLb || !current.IsMultiMaster() || !current.HasPureWorker() {
		return nil
	}
	if current.LbMode == define.ApiserverUpdater {
		err = SetupWorkersUpdater(nodes)
		if nil != err {
			return err
		}
		err = EnableWorkersUpdater(nodes)
		if nil != err {
			return err
		}
		return InstallInner(define.ApiserverUpdater, false)
	} else {
		err = SetupWorkersHaproxy(nodes)
		if nil != err {
			return err
		}
		for _, node := range nodes {
			err = action.HAProxyInitWait(current, node, 4*time.Minute)
			if nil != err {
				return err
			}
		}
		return nil
	}
}

func ChangeLoadBalance(change bool, nodes cluster.NodeList) (err error) {
	current := cluster.Current()
	if current.IsExternalLb {
		return nil
	}
	if change {
		if !current.IsMultiMaster() || !current.HasPureWorker() {
			if current.LbMode == define.ApiserverUpdater {
				_ = DeleteInner(define.ApiserverUpdater)
			}
		} else {
			return InstallLoadBalance(current.Workers)
		}
	} else {
		if nil != nodes && len(nodes) > 0 {
			return InstallLoadBalance(nodes)
		}
	}
	return nil
}

func UpgradeLoadBalance() (err error) {
	current := cluster.Current()
	if current.IsExternalLb {
		return nil
	}
	if current.IsMultiMaster() && current.HasPureWorker() {
		return InstallLoadBalance(current.Workers)
	}
	return nil
}

func SetupWorkersUpdater(nodes cluster.NodeList) (err error) {
	localConf := cluster.CurrentResource().ClusterConf
	var arr []byte
	arr, err = ShowInner(define.ApiserverStartup)
	if nil != err {
		return err
	}
	makeBalanceDir(filepath.Dir(localConf.StartupServicePath))
	err = release.WriteStartupService(arr, localConf.StartupServicePath, localConf.StartupScriptPath)
	if nil != err {
		return err
	}
	return sendNeedUpdater(nodes)
}

func EnableWorkersUpdater(nodes cluster.NodeList) (err error) {
	for _, node := range nodes {
		err = enableModuleOne(node, define.ApiserverService)
		if nil != err {
			return err
		}
	}
	return nil
}

func sendNeedUpdater(nodes cluster.NodeList) (err error) {
	localConf := cluster.CurrentResource().ClusterConf
	for _, node := range nodes {
		nodeConf := node.GetResource().ClusterConf
		err = node.CopyTo(localConf.StartupScriptPath, nodeConf.StartupScriptPath)
		if nil != err {
			return err
		}
		err = node.Chmod("+x", nodeConf.StartupScriptPath)
		if nil != err {
			return err
		}
		err = node.CopyTo(localConf.StartupServicePath, nodeConf.StartupServicePath)
		if nil != err {
			return err
		}
	}
	return nil
}

func SetupWorkersHaproxy(nodes cluster.NodeList) (err error) {
	localConf := cluster.CurrentResource().ClusterConf
	var arr []byte
	arr, err = ShowInner(define.LocalHaproxy)
	if nil != err {
		return err
	}
	makeBalanceDir(filepath.Dir(localConf.HaproxyStaticPath))
	err = onutil.WriteFile(localConf.HaproxyStaticPath, arr)
	if nil != err {
		return err
	}
	return sendNeedHaproxy(nodes)
}

func sendNeedHaproxy(nodes cluster.NodeList) (err error) {
	localConf := cluster.CurrentResource().ClusterConf
	for _, node := range nodes {
		nodeConf := node.GetResource().ClusterConf
		err = node.CopyTo(localConf.HaproxyStaticPath, nodeConf.HaproxyStaticPath)
		if nil != err {
			return err
		}
	}
	return nil
}

func makeBalanceDir(inDir string) {
	if !onutil.PathExists(inDir) {
		onutil.MkDir(inDir)
	}
}
