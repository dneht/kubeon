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
	"crypto/tls"
	cluster "github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"github.com/vbauerster/mpb/v7"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"sync"
	"time"
)

type PrepareModule struct {
	RemotePath string
	LocalPath  string
	LocalSum   string
	CopyBar    *mpb.Bar
}

func PrepareInstall(nodes cluster.NodeList, isUpgrade bool) (err error) {
	err = prepareLocal()
	if nil != err {
		return err
	}

	prog := mpb.New(
		mpb.WithWidth(90),
		mpb.WithRefreshRate(250*time.Millisecond),
	)
	sendPackage(prog, nodes, isUpgrade)
	prog.Wait()

	handlePackage(nodes, isUpgrade)
	return nil
}

func AfterUpgrade(node *cluster.Node, isBootstrap bool) (err error) {
	err = afterUpgrade(node)
	if nil != err {
		return err
	}
	if isBootstrap {
		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
		for i := 0; i < 15; i++ {
			time.Sleep(5 * time.Second)
			_, err = client.Get(node.Healthz())
			if nil == err {
				klog.V(4).Infof("Bootstrap node is ready, continue")
				return nil
			} else {
				klog.V(1).Infof("Bootstrap node is not ready, sleep 4s")
			}
		}
		return errors.New("check bootstrap health failed")
	}
	return nil
}

func prepareLocal() (err error) {
	current := cluster.Current()
	if nil == current {
		return errors.New("cluster not init")
	}

	current.UsePatch, err = release.PrepareLocal(cluster.CurrentResource())
	if nil != err {
		return err
	}
	release.AddLocalAutoCompletion()
	return nil
}

func sendPackage(prog *mpb.Progress, nodes cluster.NodeList, isUpgrade bool) {
	copyQueues := make(map[*cluster.Node][]*PrepareModule)
	klog.V(1).Info("Start copy resource to remote node")
	for _, node := range nodes {
		copyQueues[node] = doQueuedPackage(node, prog)
	}

	for node, queue := range copyQueues {
		go doTransPackage(node, queue)
	}
}

func doQueuedPackage(node *cluster.Node, prog *mpb.Progress) []*PrepareModule {
	current := cluster.Current()
	localRes := cluster.CurrentResource()
	remoteRes := node.GetResource()
	barQueue, copyQueue := make([]*mpb.Bar, 24), make([]*PrepareModule, 24)
	doAppendPackage(0, node, prog, barQueue, copyQueue, define.KubeletModule, remoteRes.KubeletPath, localRes.KubeletPath, localRes.KubeletSum)
	if current.RuntimeMode == define.ContainerdRuntime {
		doAppendPackage(1, node, prog, barQueue, copyQueue, define.ContainerdRuntime, remoteRes.ContainerdPath, localRes.ContainerdPath, localRes.ContainerdSum)
	} else {
		doAppendPackage(1, node, prog, barQueue, copyQueue, define.DockerRuntime, remoteRes.DockerPath, localRes.DockerPath, localRes.DockerSum)
	}
	doAppendPackage(2, node, prog, barQueue, copyQueue, define.NetworkPlugin, remoteRes.NetworkPath, localRes.NetworkPath, localRes.NetworkSum)
	if current.IsRealLocal() {
		doAppendPackage(3, node, prog, barQueue, copyQueue, define.ImagesPackage, remoteRes.ImagesPath, localRes.ImagesPath, localRes.ImagesSum)
		if current.IsOffline {
			doAppendPackage(4, node, prog, barQueue, copyQueue, define.OfflineModule, remoteRes.OfflinePath, localRes.OfflinePath, localRes.OfflineSum)
		}
		if current.UseKata {
			doAppendPackage(5, node, prog, barQueue, copyQueue, define.KataRuntime, remoteRes.KataPath, localRes.KataPath, localRes.KataSum)
		}
		if current.UseNvidia && node.HasNvidia {
			doAppendPackage(6, node, prog, barQueue, copyQueue, define.NvidiaRuntime, remoteRes.NvidiaPath, localRes.NvidiaPath, localRes.NvidiaSum)
		}
		switch current.IngressMode {
		case define.ContourIngress:
			{
				doAppendPackage(12, node, prog, barQueue, copyQueue, define.ContourIngress, remoteRes.ContourPath, localRes.ContourPath, localRes.ContourSum)
				break
			}
		}
	} else {
		doAppendPackage(3, node, prog, barQueue, copyQueue, define.PausePackage, remoteRes.PausePath, localRes.PausePath, localRes.PauseSum)
	}

	localScript, remoteScript := localRes.ClusterScript, remoteRes.ClusterScript
	doAppendPackage(16, node, prog, barQueue, copyQueue, define.InstallScript, remoteScript.PreparePath, localScript.PreparePath, execute.FileSum(localScript.PreparePath))
	doAppendPackage(17, node, prog, barQueue, copyQueue, define.InstallScript, remoteScript.PrepareCentosPath, localScript.PrepareCentosPath, execute.FileSum(localScript.PrepareCentosPath))
	doAppendPackage(18, node, prog, barQueue, copyQueue, define.InstallScript, remoteScript.PrepareDebianPath, localScript.PrepareDebianPath, execute.FileSum(localScript.PrepareDebianPath))
	doAppendPackage(19, node, prog, barQueue, copyQueue, define.InstallScript, remoteScript.PrepareUbuntuPath, localScript.PrepareUbuntuPath, execute.FileSum(localScript.PrepareUbuntuPath))
	doAppendPackage(20, node, prog, barQueue, copyQueue, define.InstallScript, remoteScript.DiscoverPath, localScript.DiscoverPath, execute.FileSum(localScript.DiscoverPath))
	doAppendPackage(21, node, prog, barQueue, copyQueue, define.InstallScript, remoteScript.DiscoverNvidiaPath, localScript.DiscoverNvidiaPath, execute.FileSum(localScript.DiscoverNvidiaPath))
	return copyQueue
}

func doAppendPackage(idx int, node *cluster.Node, prog *mpb.Progress, barQueue []*mpb.Bar, copyQueue []*PrepareModule, moduleName, remotePath, localPath, localSum string) {
	var prevBar *mpb.Bar
	if idx > 0 {
		for last := idx - 1; last >= 0; last-- {
			prevBar = barQueue[last]
			if nil != prevBar {
				break
			}
		}
	}
	copyBar := copyUseBar(node, prog, prevBar, moduleName, remotePath, localSum)
	if nil != copyBar {
		barQueue[idx] = copyBar
		copyQueue[idx] = &PrepareModule{remotePath, localPath, localSum, copyBar}
	}
}

func doTransPackage(node *cluster.Node, queue []*PrepareModule) {
	for _, task := range queue {
		if nil != task {
			node.CopyToWithBar(task.LocalPath, task.RemotePath, task.LocalSum, task.CopyBar)
		}
	}
}

func handlePackage(nodes cluster.NodeList, upgrade bool) {
	var wait sync.WaitGroup
	wait.Add(len(nodes))

	for _, node := range nodes {
		go doInstallPackage(&wait, node, upgrade)
	}
	wait.Wait()
}

func doInstallPackage(wait *sync.WaitGroup, node *cluster.Node, upgrade bool) {
	current := cluster.Current()
	localRes := cluster.CurrentResource()
	localConf := localRes.ClusterConf
	remoteRes := node.GetResource()
	klog.V(1).Infof("[package] start install [%s] on [%s]", define.KubeletModule, node.Addr())
	if _, err := installOnNode(node, define.KubeletModule, remoteRes.KubeletPath); nil != err {
		klog.Errorf("[package] install kubelet on [%s] failed: %v", node.Addr(), err)
		os.Exit(1)
	}
	if !upgrade {
		if err := prepareScript(node); nil != err {
			klog.Errorf("[package] prepare script on [%s] failed: %v", node.Addr(), err)
			os.Exit(1)
		}
	}
	nowRuntimeMode := cluster.Current().RuntimeMode
	klog.V(1).Infof("[package] start install [%s] on [%s]", nowRuntimeMode, node.Addr())
	if nowRuntimeMode == define.ContainerdRuntime {
		if _, err := installOnNode(node, define.ContainerdRuntime, remoteRes.ContainerdPath); nil != err {
			klog.Errorf("[package] install containerd on [%s] failed: %v", node.Addr(), err)
			os.Exit(1)
		}
	} else {
		if _, err := installOnNode(node, define.DockerRuntime, remoteRes.DockerPath); nil != err {
			klog.Errorf("[package] install docker on [%s] failed: %v", node.Addr(), err)
			os.Exit(1)
		}
	}
	if !upgrade {
		if err := enableModuleOneNow(node, nowRuntimeMode); nil != err {
			klog.Errorf("[package] enable runtime on [%s] failed: %v", node.Addr(), err)
			os.Exit(1)
		}
	}
	if err := ImportImages(current.IsRealLocal(), node); nil != err {
		klog.Errorf("[package] import container image on [%s] failed: %v", node.Addr(), err)
		os.Exit(1)
	}
	klog.V(1).Infof("[package] start install [%s] on [%s]", define.NetworkPlugin, node.Addr())
	if _, err := installOnNode(node, define.NetworkPlugin, remoteRes.NetworkPath); nil != err {
		klog.Errorf("[package] install cni on [%s] failed: %v", node.Addr(), err)
		os.Exit(1)
	}
	if current.IsRealLocal() {
		if current.UseNvidia && node.HasNvidia {
			if err := importOnNode(node, define.NvidiaRuntime, remoteRes.NvidiaPath); nil != err {
				klog.Errorf("[package] import nvidia image on [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
			setupNvidia(node, nowRuntimeMode)
		}
		if current.UseKata {
			if err := importOnNode(node, define.KataRuntime, remoteRes.KataPath); nil != err {
				klog.Errorf("[package] import kata image on [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
		}
		switch current.IngressMode {
		case define.ContourIngress:
			{
				if err := importOnNode(node, define.ContourIngress, remoteRes.ContourPath); nil != err {
					klog.Errorf("[package] import contour image on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
				break
			}
		}
	} else {
		if current.UseNvidia && node.HasNvidia {
			setupNvidia(node, nowRuntimeMode)
		}
	}
	if !upgrade {
		if err := configKubeletOne(node, localConf); nil != err {
			klog.Errorf("[package] enable kubelet on [%s] failed: %v", node.Addr(), err)
			os.Exit(1)
		}
	}
	wait.Done()
}

func prepareScript(node *cluster.Node) (err error) {
	current := cluster.Current()
	installMode := "online"
	if current.IsOffline {
		installMode = "offline"
	}
	proxyMode := current.ProxyMode
	klog.V(1).Infof("[package] start prepare install on [%s], %s, proxy=%s", node.Addr(), installMode, proxyMode)
	err = node.RunCmd("bash", node.GetResource().ScriptDir+"/prepare.sh",
		"prepare", installMode, proxyMode)
	if nil != err {
		klog.Errorf("[package] prepare install on [%s] failed", node.Addr())
		return err
	}
	if current.UseNvidia && node.HasNvidia {
		klog.V(1).Infof("[package] start discover nvidia on [%s]", node.Addr())
		err = node.RunCmd("bash", node.GetResource().ScriptDir+"/discover.sh",
			"nvidia", "no", installMode)
		if nil != err {
			klog.Errorf("[package] discover nvidia on [%s] failed", node.Addr())
			return err
		}
	}
	return nil
}

func afterUpgrade(node *cluster.Node) (err error) {
	klog.V(1).Infof("[package] start reload [%s] on [%s]", cluster.Current().RuntimeMode, node.Addr())
	if release.IsUpdateRuntime {
		if cluster.Current().RuntimeMode == define.ContainerdRuntime {
			err = restartModuleOne(node, define.ContainerdRuntime)
			if nil != err {
				return err
			}
		} else {
			err = restartModuleOne(node, define.DockerRuntime)
			if nil != err {
				return err
			}
		}
	}
	return restartModuleOne(node, define.KubeletModule)
}

func setupNvidia(node *cluster.Node, nowRuntimeMode string) {
	if nowRuntimeMode != define.ContainerdRuntime {
		klog.Errorf("[package] now runtime is not containerd, pass modify and install nvidia runtime", node.Addr())
		return
	}
	var err error
	err = node.RunCmd("sed", "-i", "-E", "\"s#BinaryName\\s+=\\s+\\\"[0-9a-zA-Z\\-]+\\\"#BinaryName = \\\"nvidia-container-runtime\\\"#g\"", "/etc/containerd/config.toml")
	if nil != err {
		klog.Errorf("[package] modify containerd config on [%s] failed: %v", node.Addr(), err)
		return
	}
	err = restartModuleOne(node, nowRuntimeMode)
	if nil != err {
		klog.Errorf("[package] restart runtime on [%s] failed: %v", node.Addr(), err)
		return
	}
}
