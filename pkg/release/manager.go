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
	"github.com/dneht/kubeon/pkg/onutil/log"
)

func InitResource(clusterVersion, runtimeMode string, isBinary, isOffline bool) *ClusterResource {
	return initResource(clusterVersion, runtimeMode, isBinary, isOffline)
}

func RemoteResource(homeDir, runtimeMode string) *ClusterRemoteResource {
	basePath := homeDir + "/.kubeon"
	return remoteResource(basePath, runtimeMode)
}

func initResource(clusterVersion, runtimeMode string, isBinary, isOffline bool) *ClusterResource {
	distPath := define.AppDistDir
	localResource := &ClusterResource{
		ImagesPath:     distPath + "/" + define.ImagesPackage + ".tar",
		ImagesSum:      onutil.GetRemoteSum(clusterVersion, define.ImagesPackage),
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
		ClusterTpl: &ClusterTplResource{
			CorednsTemplatePath: define.AppTplDir + "/" + define.CorednsPart + ".yaml.tpl",
			CalicoTemplatePath:  define.AppTplDir + "/" + define.CalicoNetwork + ".yaml.tpl",
		},
		ClusterTool: &ClusterToolResource{
			HelmPath: distPath + "/" + define.HelmTool + ".tar",
			HelmSum:  onutil.GetRemoteSum(clusterVersion, define.HelmTool),
		},
		ClusterOffline: nil,
	}
	if isBinary {
		localResource.BinaryPath = distPath + "/" + define.BinaryModule + ".tar"
		localResource.BinarySum = onutil.GetRemoteSum(clusterVersion, define.BinaryModule)
	} else {
		localResource.KubeletPath = distPath + "/" + define.KubeletModule + ".tar"
		localResource.KubeletSum = onutil.GetRemoteSum(clusterVersion, define.KubeletModule)
	}
	if isOffline {
		distOffPath := define.AppDistDir + "/" + define.OfflineModule
		localResource.ClusterOffline = &ClusterOffResource{
			OfflinePath: distOffPath + "/" + define.OfflineModule + ".tar",
			OfflineSum:  onutil.GetRemoteSum(clusterVersion, define.OfflineModule),
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
		BinaryPath:     distPath + "/" + define.BinaryModule + ".tar",
		KubeletPath:    distPath + "/" + define.KubeletModule + ".tar",
		RuntimeType:    runtimeMode,
		DockerPath:     distPath + "/" + define.DockerRuntime + ".tar",
		ContainerdPath: distPath + "/" + define.ContainerdRuntime + ".tar",
		NetworkPath:    distPath + "/" + define.NetworkPlugin + ".tar",
		OfflinePath:    distPath + "/" + define.OfflineModule + "/" + define.DockerRuntime + ".tar",
		ClusterConf: &ClusterRemoteConfResource{
			KubeletInitPath:    "/var/lib/kubelet/config.yaml",
			KubeadmInitPath:    "/etc/kubeadm.yaml",
			HaproxyStaticPath:  "/etc/kubernetes/manifests/local-haproxy.yaml",
			StartupServicePath: "/etc/systemd/system/apiserver-startup.service",
			StartupScriptPath:  "/opt/kubeon/apiserver-startup.sh",
		},
		ClusterTool: &ClusterRemoteToolResource{
			HelmPath: distPath + "/" + define.HelmTool + ".tar",
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
		log.Warnf("unpack resource failed: %s", err)
	}
	err = onutil.MvFile(localTmpDir+"/bin/kubectl", "/usr/local/bin/kubectl")
	if nil != err {
		log.Warnf("move kubectl failed: %s", err)
	}
	onutil.ChmodFile("/usr/local/bin/kubectl", 755)
	onutil.RmDir(localTmpDir)
	log.Infof("add kubectl on ")
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
			log.Warnf("set local kubectl completion failed")
		}
	}
}

func DestroyLocal() {
	onutil.RmDir(onutil.BaseDir())
}

func DelLocalAutoCompletion() {
	onutil.RmFile("/etc/profile.d/kubelet.sh")
}
