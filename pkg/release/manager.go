/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package release

import (
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/release/scriptset"
	"k8s.io/klog/v2"
)

func InitResource(clusterVersion, runtimeMode, networkMode, ingressMode string,
	isBinary, isOffline, hasNvidia, useKata, useKruise bool) *ClusterResource {
	return initResource(clusterVersion, runtimeMode, networkMode, ingressMode, isBinary, isOffline, hasNvidia, useKata, useKruise)
}

func RemoteResource(homeDir, runtimeMode, networkMode string) *ClusterRemoteResource {
	basePath := homeDir + "/.kubeon"
	return remoteResource(basePath, runtimeMode, networkMode)
}

func initResource(clusterVersion, runtimeMode, networkMode, ingressMode string,
	isBinary, isOffline, hasNvidia, useKata, useKruise bool) *ClusterResource {
	distPath := define.AppDistDir
	localResource := &ClusterResource{
		ImagesPath:  distPath + "/" + define.ImagesPackage + ".tar",
		ImagesSum:   onutil.GetRemoteSum(clusterVersion, define.ImagesPackage),
		PausePath:   distPath + "/" + define.PausePackage + ".tar",
		PauseSum:    onutil.GetRemoteSum(clusterVersion, define.PausePackage),
		RuntimeType: runtimeMode,
		NetworkType: networkMode,
		ClusterConf: &ClusterConfResource{
			KubeadmInitDir:     define.AppConfDir + "/kubeadm/",
			KubeletServicePath: define.AppConfDir + "/kubelet.service",
			KubeadmConfPath:    define.AppConfDir + "/kubeadm.conf",
			HaproxyStaticPath:  define.AppConfDir + "/haproxy/local-haproxy.yaml",
			StartupServicePath: define.AppConfDir + "/haproxy/apiserver-startup.service",
			StartupScriptPath:  define.AppConfDir + "/haproxy/apiserver-startup.sh",
		},
		ClusterScript: &ClusterScriptResource{
			PreparePath:        define.AppScriptDir + "/prepare.sh",
			PrepareCentosPath:  define.AppScriptDir + "/prepare_centos.sh",
			PrepareDebianPath:  define.AppScriptDir + "/prepare_debian.sh",
			PrepareUbuntuPath:  define.AppScriptDir + "/prepare_ubuntu.sh",
			DiscoverPath:       define.AppScriptDir + "/discover.sh",
			DiscoverNvidiaPath: define.AppScriptDir + "/discover_nvidia.sh",
			SystemVersionPath:  define.AppScriptDir + "/system_version.sh",
		},
	}
	if isBinary {
		localResource.BinaryPath = distPath + "/" + define.BinaryModule + ".tar"
		localResource.BinarySum = onutil.GetRemoteSum(clusterVersion, define.BinaryModule)
	} else {
		localResource.KubeletPath = distPath + "/" + define.KubeletModule + ".tar"
		localResource.KubeletSum = onutil.GetRemoteSum(clusterVersion, define.KubeletModule)
	}
	extVersion, _ := define.SupportComponentFull[clusterVersion]
	localResource.ContainerdPath = distPath + "/" + define.ContainerdRuntime + ".tar"
	localResource.ContainerdSum = onutil.GetRemoteSum(extVersion.Containerd, define.ContainerdRuntime)
	if runtimeMode == define.DockerRuntime {
		localResource.DockerPath = distPath + "/" + define.DockerRuntime + ".tar"
		localResource.DockerSum = onutil.GetRemoteSum(extVersion.Docker, define.DockerRuntime)
	}
	localResource.NetworkPath = distPath + "/" + define.NetworkPlugin + ".tar"
	localResource.NetworkSum = onutil.GetRemoteSum(extVersion.RealNetwork(), define.NetworkPlugin)
	if networkMode == define.CalicoNetwork {
		localResource.CalicoPath = distPath + "/" + define.CalicoNetwork + ".tar"
		localResource.CalicoSum = onutil.GetRemoteSum(extVersion.Calico, define.CalicoNetwork)
	}
	if networkMode == define.CiliumNetwork {
		localResource.CiliumPath = distPath + "/" + define.CiliumNetwork + ".tar"
		localResource.CiliumSum = onutil.GetRemoteSum(extVersion.Cilium, define.CiliumNetwork)
	}
	if ingressMode == define.ContourIngress {
		localResource.ContourPath = distPath + "/" + define.ContourIngress + ".tar"
		localResource.ContourSum = onutil.GetRemoteSum(extVersion.Contour, define.ContourIngress)
	}
	if ingressMode == define.IstioIngress {
		localResource.IstioPath = distPath + "/" + define.IstioIngress + ".tar"
		localResource.IstioSum = onutil.GetRemoteSum(extVersion.Istio, define.IstioIngress)
	}
	if hasNvidia {
		localResource.NvidiaPath = distPath + "/" + define.NvidiaRuntime + ".tar"
		localResource.NvidiaSum = onutil.GetRemoteSum(extVersion.Nvidia, define.NvidiaRuntime)
	}
	if useKata {
		localResource.KataPath = distPath + "/" + define.KataRuntime + ".tar"
		localResource.KataSum = onutil.GetRemoteSum(extVersion.Kata, define.KataRuntime)
	}
	if useKruise {
		localResource.KruisePath = distPath + "/" + define.KruisePlugin + ".tar"
		localResource.KruiseSum = onutil.GetRemoteSum(extVersion.Kruise, define.KruisePlugin)
	}
	if isOffline {
		localResource.OfflinePath = distPath + "/" + define.OfflineModule + ".tar"
		localResource.OfflineSum = onutil.GetRemoteSum(extVersion.Offline, define.OfflineModule)
	}
	writeScript(localResource.ClusterScript)
	return localResource
}

func writeScript(localScript *ClusterScriptResource) {
	onutil.MkDir(define.AppScriptDir)
	var err error
	err = onutil.WriteFile(localScript.PreparePath, []byte(scriptset.Prepare))
	if nil != err {
		panic(err)
	}
	err = onutil.WriteFile(localScript.PrepareCentosPath, []byte(scriptset.PrepareCentos))
	if nil != err {
		panic(err)
	}
	err = onutil.WriteFile(localScript.PrepareDebianPath, []byte(scriptset.PrepareDebian))
	if nil != err {
		panic(err)
	}
	err = onutil.WriteFile(localScript.PrepareUbuntuPath, []byte(scriptset.PrepareUbuntu))
	if nil != err {
		panic(err)
	}
	err = onutil.WriteFile(localScript.DiscoverPath, []byte(scriptset.Discover))
	if nil != err {
		panic(err)
	}
	err = onutil.WriteFile(localScript.DiscoverNvidiaPath, []byte(scriptset.DiscoverNvidia))
	if nil != err {
		panic(err)
	}
	err = onutil.WriteFile(localScript.SystemVersionPath, []byte(scriptset.SystemVersion))
	if nil != err {
		panic(err)
	}
}

func remoteResource(basePath, runtimeMode, networkMode string) *ClusterRemoteResource {
	distPath := basePath + "/dist"
	scriptPath := basePath + "/script"
	return &ClusterRemoteResource{
		BaseDir:        basePath,
		ConfDir:        basePath + "/conf",
		TplDir:         basePath + "/tpl",
		ScriptDir:      scriptPath,
		PatchDir:       basePath + "/patch",
		DistDir:        distPath,
		TmpDir:         basePath + "/tmp",
		ImagesPath:     distPath + "/" + define.ImagesPackage + ".tar",
		PausePath:      distPath + "/" + define.PausePackage + ".tar",
		BinaryPath:     distPath + "/" + define.BinaryModule + ".tar",
		KubeletPath:    distPath + "/" + define.KubeletModule + ".tar",
		OfflinePath:    distPath + "/" + define.OfflineModule + ".tar",
		RuntimeType:    runtimeMode,
		DockerPath:     distPath + "/" + define.DockerRuntime + ".tar",
		ContainerdPath: distPath + "/" + define.ContainerdRuntime + ".tar",
		NetworkType:    networkMode,
		NetworkPath:    distPath + "/" + define.NetworkPlugin + ".tar",
		CalicoPath:     distPath + "/" + define.CalicoNetwork + ".tar",
		CiliumPath:     distPath + "/" + define.CiliumNetwork + ".tar",
		ContourPath:    distPath + "/" + define.ContourIngress + ".tar",
		IstioPath:      distPath + "/" + define.IstioIngress + ".tar",
		NvidiaPath:     distPath + "/" + define.NvidiaRuntime + ".tar",
		KataPath:       distPath + "/" + define.KataRuntime + ".tar",
		KruisePath:     distPath + "/" + define.KruisePlugin + ".tar",
		ClusterConf: &ClusterRemoteConfResource{
			KubeletInitPath:    "/var/lib/kubelet/config.yaml",
			KubeadmInitPath:    "/etc/kubeadm.yaml",
			HaproxyStaticPath:  "/etc/kubernetes/manifests/local-haproxy.yaml",
			StartupServicePath: "/etc/systemd/system/apiserver-startup.service",
			StartupScriptPath:  "/opt/kubeon/apiserver-startup.sh",
		},
		ClusterScript: &ClusterRemoteScriptResource{
			PreparePath:        scriptPath + "/prepare.sh",
			PrepareCentosPath:  scriptPath + "/prepare_centos.sh",
			PrepareDebianPath:  scriptPath + "/prepare_debian.sh",
			PrepareUbuntuPath:  scriptPath + "/prepare_ubuntu.sh",
			DiscoverPath:       scriptPath + "/discover.sh",
			DiscoverNvidiaPath: scriptPath + "/discover_nvidia.sh",
			SystemVersionPath:  scriptPath + "/system_version.sh",
		},
	}
}

func PrepareLocal(resource *ClusterResource) (use bool, err error) {
	createKubeonNeedDirs()

	localTmpDir := define.AppTmpDir + "/local"
	onutil.MkDir(localTmpDir)
	err = execute.UnpackTar(resource.KubeletPath, localTmpDir)
	if nil != err {
		return false, err
	}
	err = onutil.MvDir(localTmpDir+"/tpl", define.AppTplDir)
	if nil != err {
		return false, err
	}
	err = onutil.MvFile(localTmpDir+"/bin/kubectl", "/usr/local/bin/kubectl")
	if nil != err {
		return false, err
	}
	onutil.ChmodFile("/usr/local/bin/kubectl", 755)
	err = onutil.MvFile(localTmpDir+"/bin/kubeadm", "/usr/local/bin/kubeadm")
	if nil != err {
		return false, err
	}
	onutil.ChmodFile("/usr/local/bin/kubeadm", 755)
	use = onutil.IsEmptyDir(localTmpDir + "/patch")
	onutil.RmDir(localTmpDir)

	cniTmpDir := define.AppTmpDir + "/cni"
	onutil.MkDir(cniTmpDir)
	err = execute.UnpackTar(resource.NetworkPath, cniTmpDir)
	if nil != err {
		return false, err
	}
	err = execute.NewLocalCmd("bash", cniTmpDir+"/install.sh").Run()
	if nil != err {
		return false, err
	}
	onutil.RmDir(cniTmpDir)
	return use, nil
}

func ReinstallLocal(resource *ClusterResource) {
	createKubeonNeedDirs()

	localTmpDir := define.AppTmpDir + "/local"
	onutil.MkDir(localTmpDir)
	err := execute.UnpackTar(resource.KubeletPath, localTmpDir)
	if nil != err {
		klog.Warningf("Unpack resource failed: %s", err)
	}
	err = onutil.MvFile(localTmpDir+"/bin/kubectl", "/usr/local/bin/kubectl")
	if nil != err {
		klog.Warningf("Move kubectl failed: %s", err)
	}
	onutil.ChmodFile("/usr/local/bin/kubectl", 755)
	onutil.RmDir(localTmpDir)
	klog.V(1).Infof("Local host add kubectl on ")
	AddLocalAutoCompletion()
}

func createKubeonNeedDirs() {
	onutil.MkDir(define.AppTmpDir)
	onutil.MkDir(define.AppConfDir)
	onutil.MkDir(define.AppPatchDir)
	onutil.MkDir(define.AppTplDir)
}

func AddLocalAutoCompletion() {
	onutil.MkDir("/etc/profile.d")
	kubeAutoPath := "/etc/profile.d/kubelet.sh"
	if !onutil.PathExists(kubeAutoPath) {
		if err := execute.NewLocalCmd("sh", "-c", "echo 'source <(kubectl completion bash)' >"+kubeAutoPath).Run(); nil != err {
			klog.Warningf("Set local kubectl completion failed")
		}
		if err := execute.NewLocalCmd("sh", "-c", "echo 'source <(crictl completion bash)' >>"+kubeAutoPath).Run(); nil != err {
			klog.Warningf("Set local crictl completion failed")
		}
	}
	kubeonAutoPath := "/etc/profile.d/kubeon.sh"
	if !onutil.PathExists(kubeonAutoPath) {
		if err := execute.NewLocalCmd("sh", "-c", "echo 'source <(kubeon completion bash)' >"+kubeonAutoPath).Run(); nil != err {
			klog.Warningf("Set local kubeon completion failed")
		}
	}
}

func DestroyLocal() {
	onutil.RmFile("/usr/local/bin/kubectl")
	onutil.RmFile("/usr/local/bin/kubeadm")
	onutil.RmDir(onutil.BaseDir())
	_ = execute.NewLocalCmd("bash", define.AppScriptDir+"/uninstall_cni.sh").Run()
}

func DelLocalAutoCompletion() {
	onutil.RmFile("/etc/profile.d/kubelet.sh")
	onutil.RmFile("/etc/profile.d/kubeon.sh")
}
