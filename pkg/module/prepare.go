/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package module

import (
	"crypto/tls"
	cluster "github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
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

	lch := make(chan struct{}, 6)
	for node, queue := range copyQueues {
		lch <- struct{}{}
		go doTransPackage(node, queue, lch)
	}
}

func doQueuedPackage(node *cluster.Node, prog *mpb.Progress) []*PrepareModule {
	current := cluster.Current()
	localRes := cluster.CurrentResource()
	remoteRes := node.GetResource()
	barQueue, copyQueue := make([]*mpb.Bar, 24), make([]*PrepareModule, 24)
	doAppendPackage(0, node, prog, barQueue, copyQueue, define.KubeletModule, remoteRes.KubeletPath, localRes.KubeletPath, localRes.KubeletSum)
	doAppendPackage(1, node, prog, barQueue, copyQueue, define.ScriptsModule, remoteRes.ScriptsPath, localRes.ScriptsPath, localRes.ScriptsSum)
	if current.RuntimeMode == define.ContainerdRuntime {
		doAppendPackage(4, node, prog, barQueue, copyQueue, define.ContainerdRuntime, remoteRes.ContainerdPath, localRes.ContainerdPath, localRes.ContainerdSum)
	} else {
		doAppendPackage(4, node, prog, barQueue, copyQueue, define.DockerRuntime, remoteRes.DockerPath, localRes.DockerPath, localRes.DockerSum)
	}
	doAppendPackage(5, node, prog, barQueue, copyQueue, define.NetworkPlugin, remoteRes.NetworkPath, localRes.NetworkPath, localRes.NetworkSum)
	if current.IsRealLocal() {
		doAppendPackage(6, node, prog, barQueue, copyQueue, define.ImagesPackage, remoteRes.ImagesPath, localRes.ImagesPath, localRes.ImagesSum)
		if current.IsOffline {
			doAppendPackage(7, node, prog, barQueue, copyQueue, define.OfflineModule, remoteRes.OfflinePath, localRes.OfflinePath, localRes.OfflineSum)
		}
		if current.UseNvidia && node.HasNvidia {
			doAppendPackage(8, node, prog, barQueue, copyQueue, define.NvidiaRuntime, remoteRes.NvidiaPath, localRes.NvidiaPath, localRes.NvidiaSum)
		}
		if current.UseKata {
			doAppendPackage(9, node, prog, barQueue, copyQueue, define.KataRuntime, remoteRes.KataPath, localRes.KataPath, localRes.KataSum)
		}
		switch current.NetworkMode {
		case define.CalicoNetwork:
			{
				doAppendPackage(14, node, prog, barQueue, copyQueue, define.CalicoNetwork, remoteRes.CalicoPath, localRes.CalicoPath, localRes.CalicoSum)
				break
			}
		case define.CiliumNetwork:
			{
				doAppendPackage(14, node, prog, barQueue, copyQueue, define.CiliumNetwork, remoteRes.CiliumPath, localRes.CiliumPath, localRes.CiliumSum)
				break
			}
		}
		switch current.IngressMode {
		case define.ContourIngress:
			{
				doAppendPackage(16, node, prog, barQueue, copyQueue, define.ContourIngress, remoteRes.ContourPath, localRes.ContourPath, localRes.ContourSum)
				break
			}
		case define.IstioIngress:
			{
				doAppendPackage(16, node, prog, barQueue, copyQueue, define.IstioIngress, remoteRes.IstioPath, localRes.IstioPath, localRes.IstioSum)
				break
			}
		}
		if current.UseKruise {
			doAppendPackage(18, node, prog, barQueue, copyQueue, define.KruisePlugin, remoteRes.KruisePath, localRes.KruisePath, localRes.KruiseSum)
		}
	} else {
		doAppendPackage(6, node, prog, barQueue, copyQueue, define.PausePackage, remoteRes.PausePath, localRes.PausePath, localRes.PauseSum)
	}
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

func doTransPackage(node *cluster.Node, queue []*PrepareModule, lch chan struct{}) {
	for _, task := range queue {
		if nil != task {
			node.CopyToWithBar(task.LocalPath, task.RemotePath, task.LocalSum, task.CopyBar)
		}
	}
	<-lch
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
	if _, err := installOnNode(node, define.ScriptsModule, remoteRes.ScriptsPath); nil != err {
		klog.Errorf("[package] install scripts on [%s] failed: %v", node.Addr(), err)
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
		pauseImageAddr := current.GetPauseImageAddr()
		if "" != pauseImageAddr {
			err := node.RunCmd("sed", "-i", "-E", "\"s#sandbox_image\\s+=\\s+\\\"[0-9a-zA-Z\\/+-.:]+\\\"#sandbox_image = \\\""+pauseImageAddr+"\\\"#g\"", "/etc/containerd/config.toml")
			if nil != err {
				klog.Errorf("[package] modify containerd infra image on [%s] failed: %v", node.Addr(), err)
			}
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
		switch current.NetworkMode {
		case define.CalicoNetwork:
			{
				if err := importOnNode(node, define.CalicoNetwork, remoteRes.CalicoPath); nil != err {
					klog.Errorf("[package] import calico image on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
				break
			}
		case define.CiliumNetwork:
			{
				if err := importOnNode(node, define.CiliumNetwork, remoteRes.CiliumPath); nil != err {
					klog.Errorf("[package] import cilium image on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
				break
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
		case define.IstioIngress:
			{
				if err := importOnNode(node, define.IstioIngress, remoteRes.IstioPath); nil != err {
					klog.Errorf("[package] import istio image on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
				break
			}
		}
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
		if current.UseKruise {
			if err := importOnNode(node, define.KruisePlugin, remoteRes.KruisePath); nil != err {
				klog.Errorf("[package] import kruise image on [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
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
