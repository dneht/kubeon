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

package module

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/release"
)

func ConfigKubelet(nodeSelector string) (err error) {
	getNodes := selectNodes(nodeSelector)
	localConf := cluster.CurrentResource().ClusterConf

	for _, node := range getNodes {
		err = configKubeletOne(node, localConf)
		if nil != err {
			return err
		}
	}
	return nil
}

func configKubeletOne(node *cluster.Node, localConf *release.ClusterConfResource) (err error) {
	return enableModuleOne(node, define.KubeletModule)
}

func EnableModule(moduleName, nodeSelector string) (err error) {
	getNodes := selectNodes(nodeSelector)

	for _, node := range getNodes {
		err = enableModuleOneNow(node, moduleName)
		if nil != err {
			return err
		}
	}
	return nil
}

func enableModuleOne(node *cluster.Node, moduleName string) (err error) {
	return node.RunCmd("systemctl", "enable", moduleName)
}

func enableModuleOneNow(node *cluster.Node, moduleName string) (err error) {
	return node.RunCmd("systemctl", "enable", moduleName, "--now")
}

func RestartModule(moduleName, nodeSelector string) (err error) {
	getNodes := selectNodes(nodeSelector)

	for _, node := range getNodes {
		err = restartModuleOne(node, moduleName)
		if nil != err {
			return err
		}
	}
	return nil
}

func restartModuleOne(node *cluster.Node, moduleName string) (err error) {
	return node.RunCmd("systemctl", "daemon-reload", "&&", "systemctl", "restart", moduleName)
}
