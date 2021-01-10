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
	ImagesPath     string
	ImagesSum      string
	BinaryPath     string
	BinarySum      string
	KubeletPath    string
	KubeletSum     string
	RuntimeType    string
	DockerPath     string
	DockerSum      string
	ContainerdPath string
	ContainerdSum  string
	NetworkPath    string
	NetworkSum     string
	ClusterConf    *ClusterConfResource
	ClusterTpl     *ClusterTplResource
	ClusterTool    *ClusterToolResource
	ClusterOffline *ClusterOffResource
}

type ClusterConfResource struct {
	KubeadmInitDir     string
	KubeletServicePath string
	KubeadmConfPath    string
	HaproxyStaticPath  string
	StartupServicePath string
	StartupScriptPath  string
}

type ClusterTplResource struct {
	CorednsTemplatePath string
	CalicoTemplatePath  string
}

type ClusterOffResource struct {
	OfflinePath string
	OfflineSum  string
}

type ClusterToolResource struct {
	HelmPath string
	HelmSum  string
}

type ClusterRemoteResource struct {
	BaseDir        string
	ConfDir        string
	TplDir         string
	ScriptDir      string
	PatchDir       string
	DistDir        string
	TmpDir         string
	ImagesPath     string
	BinaryPath     string
	KubeletPath    string
	RuntimeType    string
	DockerPath     string
	ContainerdPath string
	NetworkPath    string
	OfflinePath    string
	ClusterConf    *ClusterRemoteConfResource
	ClusterTool    *ClusterRemoteToolResource
}

type ClusterRemoteConfResource struct {
	KubeletInitPath    string
	KubeadmInitPath    string
	HaproxyStaticPath  string
	StartupServicePath string
	StartupScriptPath  string
}

type ClusterRemoteToolResource struct {
	HelmPath string
}
