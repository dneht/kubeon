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

package release

import (
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/mitchellh/mapstructure"
	"k8s.io/klog/v2"
)

func InitResource(clusterVersion, runtimeMode string, isBinary, isOffline, hasNvidia, useKata bool, ingressMode string) *ClusterResource {
	return initResource(clusterVersion, runtimeMode, isBinary, isOffline, hasNvidia, useKata, ingressMode)
}

func RemoteResource(homeDir, runtimeMode string) *ClusterRemoteResource {
	basePath := homeDir + "/.kubeon"
	return remoteResource(basePath, runtimeMode)
}

func initResource(clusterVersion, runtimeMode string, isBinary, isOffline,
	hasNvidia, useKata bool, ingressMode string) *ClusterResource {
	distPath := define.AppDistDir
	localResource := &ClusterResource{
		ImagesPath:     distPath + "/" + define.ImagesPackage + ".tar",
		ImagesSum:      onutil.GetRemoteSum(clusterVersion, define.ImagesPackage),
		PausePath:      distPath + "/" + define.PausePackage + ".tar",
		PauseSum:       onutil.GetRemoteSum(clusterVersion, define.PausePackage),
		RuntimeType:    runtimeMode,
		DockerPath:     distPath + "/" + define.DockerRuntime + ".tar",
		DockerSum:      onutil.GetRemoteSum(clusterVersion, define.DockerRuntime),
		ContainerdPath: distPath + "/" + define.ContainerdRuntime + ".tar",
		ContainerdSum:  onutil.GetRemoteSum(clusterVersion, define.ContainerdRuntime),
		NetworkPath:    distPath + "/" + define.NetworkPlugin + ".tar",
		NetworkSum:     onutil.GetRemoteSum(clusterVersion, define.NetworkPlugin),
		ClusterConf: &ClusterConfResource{
			KubeadmInitDir:     define.AppConfDir + "/kubeadm/",
			KubeletServicePath: define.AppConfDir + "/kubelet.service",
			KubeadmConfPath:    define.AppConfDir + "/kubeadm.conf",
			HaproxyStaticPath:  define.AppConfDir + "/haproxy/local-haproxy.yaml",
			StartupServicePath: define.AppConfDir + "/haproxy/apiserver-startup.service",
			StartupScriptPath:  define.AppConfDir + "/haproxy/apiserver-startup.sh",
		},
	}
	if isBinary {
		localResource.BinaryPath = distPath + "/" + define.BinaryModule + ".tar"
		localResource.BinarySum = onutil.GetRemoteSum(clusterVersion, define.BinaryModule)
	} else {
		localResource.KubeletPath = distPath + "/" + define.KubeletModule + ".tar"
		localResource.KubeletSum = onutil.GetRemoteSum(clusterVersion, define.KubeletModule)
	}
	extVersion, ok := define.SupportComponentFull[clusterVersion]
	if ok {
		if hasNvidia {
			localResource.NvidiaPath = distPath + "/" + define.NvidiaRuntime + ".tar"
			localResource.NvidiaSum = onutil.GetRemoteSum(extVersion.Nvidia, define.NvidiaRuntime)
		}
		if useKata {
			localResource.KataPath = distPath + "/" + define.KataRuntime + ".tar"
			localResource.KataSum = onutil.GetRemoteSum(extVersion.Kata, define.KataRuntime)
		}
		if ingressMode == define.ContourIngress {
			localResource.ContourPath = distPath + "/" + define.ContourIngress + ".tar"
			localResource.ContourSum = onutil.GetRemoteSum(extVersion.Contour, define.ContourIngress)
		}
	}
	if isOffline {
		localResource.OfflinePath = distPath + "/" + define.OfflineModule + ".tar"
		localResource.OfflineSum = onutil.GetRemoteSum(clusterVersion, define.OfflineModule)
	}
	installVer := define.SupportComponentFull[clusterVersion]
	if nil != installVer {
		installVerMap := map[string]string{}
		err := mapstructure.Decode(installVer, &installVerMap)
		if nil != err {
			localResource.InstallVersion = &installVerMap
		}
	}
	return localResource
}

func remoteResource(basePath, runtimeMode string) *ClusterRemoteResource {
	distPath := basePath + "/dist"
	return &ClusterRemoteResource{
		BaseDir:        basePath,
		ConfDir:        basePath + "/conf",
		TplDir:         basePath + "/tpl",
		ScriptDir:      basePath + "/script",
		PatchDir:       basePath + "/patch",
		DistDir:        distPath,
		TmpDir:         basePath + "/tmp",
		ImagesPath:     distPath + "/" + define.ImagesPackage + ".tar",
		PausePath:      distPath + "/" + define.PausePackage + ".tar",
		BinaryPath:     distPath + "/" + define.BinaryModule + ".tar",
		KubeletPath:    distPath + "/" + define.KubeletModule + ".tar",
		RuntimeType:    runtimeMode,
		DockerPath:     distPath + "/" + define.DockerRuntime + ".tar",
		ContainerdPath: distPath + "/" + define.ContainerdRuntime + ".tar",
		NvidiaPath:     distPath + "/" + define.NvidiaRuntime + ".tar",
		KataPath:       distPath + "/" + define.KataRuntime + ".tar",
		NetworkPath:    distPath + "/" + define.NetworkPlugin + ".tar",
		ContourPath:    distPath + "/" + define.ContourIngress + ".tar",
		OfflinePath:    distPath + "/" + define.OfflineModule + ".tar",
		ClusterConf: &ClusterRemoteConfResource{
			KubeletInitPath:    "/var/lib/kubelet/config.yaml",
			KubeadmInitPath:    "/etc/kubeadm.yaml",
			HaproxyStaticPath:  "/etc/kubernetes/manifests/local-haproxy.yaml",
			StartupServicePath: "/etc/systemd/system/apiserver-startup.service",
			StartupScriptPath:  "/opt/kubeon/apiserver-startup.sh",
		},
	}
}

func PrepareLocal(resource *ClusterResource) (use bool, err error) {
	localTmpDir := define.AppTmpDir + "/local"
	mkNeedDir(localTmpDir)
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
	return use, nil
}

func ReinstallLocal(resource *ClusterResource) {
	localTmpDir := define.AppTmpDir + "/local"
	mkNeedDir(localTmpDir)
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

func mkNeedDir(localTmpDir string) {
	onutil.MkDir(localTmpDir)
	onutil.MkDir(define.AppTmpDir)
	onutil.MkDir(define.AppConfDir)
	onutil.MkDir(define.AppPatchDir)
	onutil.MkDir(define.AppTplDir)
}

func AddLocalAutoCompletion() {
	onutil.MkDir("/etc/profile.d")
	localAutoPath := "/etc/profile.d/kubelet.sh"
	if !onutil.PathExists(localAutoPath) {
		err := execute.NewLocalCmd("sh", "-c", "echo 'source <(kubectl completion bash)' >>"+localAutoPath).Run()
		if nil != err {
			klog.Warningf("Set local kubectl completion failed")
		}
	}
}

func DestroyLocal() {
	onutil.RmFile("/usr/local/bin/kubectl")
	onutil.RmFile("/usr/local/bin/kubeadm")
	onutil.RmDir(onutil.BaseDir())
}

func DelLocalAutoCompletion() {
	onutil.RmFile("/etc/profile.d/kubelet.sh")
}
