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

package release

type ClusterResource struct {
	ImagesPath     string               `json:"imagesPath"`
	ImagesSum      string               `json:"imagesSum"`
	BinaryPath     string               `json:"binaryPath"`
	BinarySum      string               `json:"binarySum"`
	KubeletPath    string               `json:"kubeletPath"`
	KubeletSum     string               `json:"kubeletSum"`
	RuntimeType    string               `json:"runtimeType"`
	DockerPath     string               `json:"dockerPath"`
	DockerSum      string               `json:"dockerSum"`
	ContainerdPath string               `json:"containerdPath"`
	ContainerdSum  string               `json:"containerdSum"`
	NetworkPath    string               `json:"networkPath"`
	NetworkSum     string               `json:"networkSum"`
	ClusterConf    *ClusterConfResource `json:"clusterConf"`
	ClusterTpl     *ClusterTplResource  `json:"clusterTpl"`
	ClusterOffline *ClusterOffResource  `json:"clusterOffline"`
}

type ClusterConfResource struct {
	KubeadmInitDir     string `json:"kubeadmInitDir"`
	KubeletServicePath string `json:"kubeletServicePath"`
	KubeadmConfPath    string `json:"kubeadmConfPath"`
	HaproxyStaticPath  string `json:"haproxyStaticPath"`
	StartupServicePath string `json:"startupServicePath"`
	StartupScriptPath  string `json:"startupScriptPath"`
}

type ClusterTplResource struct {
	CorednsTemplatePath string `json:"corednsTemplatePath"`
	CalicoTemplatePath  string `json:"calicoTemplatePath"`
}

type ClusterOffResource struct {
	OfflinePath string `json:"offlinePath"`
	OfflineSum  string `json:"offlineSum"`
}

type ClusterRemoteResource struct {
	BaseDir        string                     `json:"baseDir"`
	ConfDir        string                     `json:"confDir"`
	TplDir         string                     `json:"tplDir"`
	ScriptDir      string                     `json:"scriptDir"`
	PatchDir       string                     `json:"patchDir"`
	DistDir        string                     `json:"distDir"`
	TmpDir         string                     `json:"tmpDir"`
	ImagesPath     string                     `json:"imagesPath"`
	BinaryPath     string                     `json:"binaryPath"`
	KubeletPath    string                     `json:"kubeletPath"`
	RuntimeType    string                     `json:"runtimeType"`
	DockerPath     string                     `json:"dockerPath"`
	ContainerdPath string                     `json:"containerdPath"`
	NetworkPath    string                     `json:"networkPath"`
	OfflinePath    string                     `json:"offlinePath"`
	ClusterConf    *ClusterRemoteConfResource `json:"clusterConf"`
}

type ClusterRemoteConfResource struct {
	KubeletInitPath    string `json:"kubeletInitPath"`
	KubeadmInitPath    string `json:"kubeadmInitPath"`
	HaproxyStaticPath  string `json:"haproxyStaticPath"`
	StartupServicePath string `json:"startupServicePath"`
	StartupScriptPath  string `json:"startupScriptPath"`
}
