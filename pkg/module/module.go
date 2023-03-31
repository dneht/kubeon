/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package module

import (
	cluster "github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"k8s.io/klog/v2"
	"os"
	"time"
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

	progGroup := mpb.New(
		mpb.WithWidth(90),
		mpb.WithRefreshRate(250*time.Millisecond),
	)
	getNodes := selectNodes(nodeSelector)
	for _, node := range getNodes {
		remote := remotePath(node, moduleName)
		bar := copyUseBar(node, progGroup, nil, moduleName, remote, localSum)
		if nil != bar {
			go node.CopyToWithBar(localPath, remote, localSum, bar)
		}
	}
	progGroup.Wait()

	for _, node := range getNodes {
		if res, ierr := installOnNode(node, moduleName, remotePath(node, moduleName)); nil != ierr {
			return ierr
		} else {
			if res {
				eerr := enableModuleOneNow(node, moduleName)
				if nil != eerr {
					klog.Warningf("Enable systemd failed: %v", err)
				}
			}
		}
	}
	return nil
}

func remotePath(node *cluster.Node, module string) string {
	var nodePath string
	nodeResource := node.GetResource()
	switch module {
	case define.DockerRuntime:
		nodePath = nodeResource.DockerPath
		break
	case define.ContainerdRuntime:
		nodePath = nodeResource.ContainerdPath
		break
	case define.NetworkPlugin:
		nodePath = nodeResource.NetworkPath
		break
	default:
		nodePath = nodeResource.DistDir + "/" + module + ".tar"
	}
	return nodePath
}

func copyUseBar(node *cluster.Node, prog *mpb.Progress, bar *mpb.Bar, module, remotePath, localSum string) *mpb.Bar {
	remoteSum := node.FileSum(remotePath)
	if localSum != remoteSum {
		if nil == bar {
			return prog.New(0,
				mpb.BarStyle().Rbound("|"),
				mpb.PrependDecorators(
					decor.Name("copy "+module+" to "+node.Hostname, decor.WC{W: 32, C: decor.DidentRight}),
					decor.CountersKibiByte("% .2f / % .2f"),
				),
				mpb.AppendDecorators(
					decor.Percentage(decor.WC{W: 5}),
					decor.Name(" ] "),
				),
			)
		} else {
			return prog.AddBar(0,
				mpb.BarQueueAfter(bar, false),
				mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					decor.Name("copy "+module+" to "+node.Hostname, decor.WC{W: 32, C: decor.DidentRight}),
					decor.CountersKibiByte("% .2f / % .2f"),
				),
				mpb.AppendDecorators(
					decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
					decor.OnComplete(decor.Name(" ] "), ""),
				),
			)
		}
	} else {
		return nil
	}
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
	err = importImages(node, moduleName, remoteTmpDir+"/image.tar")
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
