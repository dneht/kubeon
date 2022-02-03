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
	cluster "github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"k8s.io/klog/v2"
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
		localPath = define.AppDistDir + "/" + moduleName + ".tar"
		if onutil.PathExists(localPath) {
			localSum = execute.FileSum(localPath)
		} else {
			klog.Warningf("Not support module[%s]", moduleName)
			return nil
		}
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
		default:
			remotePath = remoteResource.DistDir + "/" + moduleName + ".tar"
		}
		err = copyToNode(node, remotePath, localPath, localSum)
		if nil != err {
			return err
		}
		var sd bool
		sd, err = installOnNode(node, moduleName, remotePath)
		if nil != err {
			return err
		}
		if sd {
			err = enableModuleOneNow(node, moduleName)
			if nil != err {
				klog.Warningf("Enable systemd failed: %v", err)
			}
		}
	}
	return nil
}

func copyToNode(node *cluster.Node, remotePath, localPath, localSum string) (err error) {
	remoteSum := node.FileSum(remotePath)
	klog.V(4).Infof("Get local[%s] sum %s, remote[%s] sum %s", localPath, localSum, remotePath, remoteSum)
	if localSum != remoteSum {
		err = node.CopyToWithSum(localPath, remotePath, localSum)
		if nil != err {
			return err
		}
	}
	return nil
}

func unzipOnNode(node *cluster.Node, moduleName, remotePath string) (err error) {
	remoteResource := node.GetResource()
	remoteTmpDir := remoteResource.TmpDir + "/" + moduleName
	err = node.RunCmd("mkdir", "-p", remoteTmpDir, "&&", "tar", "xf", remotePath, "-C", remoteTmpDir)
	if nil != err {
		return err
	}
	return nil
}

func installOnNode(node *cluster.Node, moduleName, remotePath string) (sd bool, err error) {
	remoteResource := node.GetResource()
	scriptDir := remoteResource.ScriptDir
	remoteTmpDir := remoteResource.TmpDir + "/" + moduleName
	err = unzipOnNode(node, moduleName, remotePath)
	if nil != err {
		return false, err
	}
	sd = false
	result, err := node.Command("if", "test", "-f", remoteTmpDir+"/"+moduleName+".service", ";then",
		"echo", "yes;", "fi").RunAndResult()
	if nil == err && "yes" == result {
		sd = true
	}
	current := cluster.Current()
	err = node.RunCmd("bash", remoteTmpDir+"/install.sh", current.Version.String(), current.RuntimeMode)
	if nil != err {
		return sd, err
	}
	err = node.RunCmd("if", "test", "-f", remoteTmpDir+"/uninstall.sh", ";then",
		"mv", "-f", remoteTmpDir+"/uninstall.sh", scriptDir+"/uninstall_"+moduleName+".sh", ";fi")
	if nil != err {
		_ = node.RunCmd("bash", scriptDir+"/uninstall_"+moduleName+".sh")
		return sd, err
	}
	_ = node.RunCmd("rm", "-rf", remoteTmpDir)
	return sd, nil
}

func importOnNode(node *cluster.Node, moduleName, remotePath string) (err error) {
	remoteResource := node.GetResource()
	remoteTmpDir := remoteResource.TmpDir + "/" + moduleName
	err = unzipOnNode(node, moduleName, remotePath)
	if nil != err {
		return err
	}
	err = importImage(node, remoteTmpDir+"/image.tar")
	_ = node.RunCmd("rm", "-rf", remoteTmpDir)
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
	remoteResource := node.GetResource()
	scriptDir := remoteResource.ScriptDir
	unScript := scriptDir + "/uninstall_" + moduleName + ".sh"
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
		klog.Error("Cluster not exist, please check context")
		os.Exit(1)
	}

	getNodes, err := cluster.SelectNodes(nodeSelector)
	if nil != err {
		klog.Errorf("Cluster node select[%s] failed, please check context", nodeSelector)
		os.Exit(1)
	}
	if len(getNodes) == 0 {
		klog.Warningf("Cluster node select[%s] is empty, do noting", nodeSelector)
	}
	return getNodes
}
