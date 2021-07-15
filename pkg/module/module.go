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
	"github.com/dneht/kubeon/pkg/onutil/log"
	"os"
)

func InstallSelect(moduleName, nodeSelector string) (err error) {
	var localPath string
	var localSum string
	currentResource := cluster.CurrentResource()
	switch moduleName {
	case define.DockerRuntime:
		localPath = currentResource.DockerPath
		localSum = currentResource.DockerSum
		break
	case define.ContainerdRuntime:
		localPath = currentResource.ContainerdPath
		localSum = currentResource.ContainerdSum
		break
	case define.NetworkPlugin:
		localPath = currentResource.NetworkPath
		localSum = currentResource.NetworkSum
		break
	default:
		log.Warnf("not support module[%s]", moduleName)
		return nil
	}

	getNodes := selectNodes(nodeSelector)
	for _, node := range getNodes {
		var remotePath string
		remoteResource := node.GetResource()
		switch moduleName {
		case define.DockerRuntime:
			remotePath = remoteResource.DockerPath
			break
		case define.ContainerdRuntime:
			remotePath = remoteResource.ContainerdPath
			break
		case define.NetworkPlugin:
			remotePath = remoteResource.NetworkPath
			break
		}
		err = copyToNode(node, remotePath, localPath, localSum)
		if nil != err {
			return err
		}
		err = installOnNode(node, moduleName, remotePath)
		if nil != err {
			return err
		}
	}
	return nil
}

func copyToNode(node *cluster.Node, remotePath, localPath, localSum string) (err error) {
	remoteSum := node.FileSum(remotePath)
	log.Debugf("get local[%s] sum %s, remote[%s] sum %s", localPath, localSum, remotePath, remoteSum)
	if localSum != remoteSum {
		err = node.CopyToWithSum(localPath, remotePath, localSum)
		if nil != err {
			return err
		}
	}
	return nil
}

func installOnNode(node *cluster.Node, moduleName, remotePath string) (err error) {
	remoteResource := node.GetResource()
	remoteTmpDir := remoteResource.TmpDir + "/" + moduleName
	err = node.RunCmd("mkdir", "-p", remoteTmpDir,
		"&&", "tar", "xf", remotePath, "-C", remoteTmpDir,
		"&&", "bash", remoteTmpDir+"/install.sh")
	if nil != err {
		return err
	}
	return nil
}

func UninstallSelect(moduleName, nodeSelector string) error {
	getNodes := selectNodes(nodeSelector)
	var err error
	for _, node := range getNodes {
		err = uninstallOnNode(node, moduleName)
		if nil != err {
			return err
		}
	}
	return nil
}

func uninstallOnNode(node *cluster.Node, moduleName string) (err error) {
	unScript := node.GetResource().ScriptDir + "/uninstall_" + moduleName + ".sh"
	exist := node.FileExist(unScript)
	if exist {
		err = node.RunCmd("bash", unScript)
		if nil != err {
			return err
		}
	}
	return nil
}

func selectNodes(nodeSelector string) cluster.NodeList {
	if nil == cluster.Current() {
		log.Error("cluster not exist, please check context")
		os.Exit(1)
	}

	getNodes, err := cluster.SelectNodes(nodeSelector)
	if nil != err {
		log.Errorf("cluster node select[%s] failed, please check context", nodeSelector)
		os.Exit(1)
	}
	if len(getNodes) == 0 {
		log.Warnf("cluster node select[%s] is empty, do noting", nodeSelector)
	}
	return getNodes
}
