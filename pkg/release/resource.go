/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package release

type ClusterResource struct {
	ImagesPath     string                 `json:"imagesPath"`
	ImagesSum      string                 `json:"imagesSum"`
	PausePath      string                 `json:"pausePath"`
	PauseSum       string                 `json:"pauseSum"`
	BinaryPath     string                 `json:"binaryPath,omitempty"`
	BinarySum      string                 `json:"binarySum,omitempty"`
	KubeletPath    string                 `json:"kubeletPath,omitempty"`
	KubeletSum     string                 `json:"kubeletSum,omitempty"`
	OfflinePath    string                 `json:"offlinePath,omitempty"`
	OfflineSum     string                 `json:"offlineSum,omitempty"`
	RuntimeType    string                 `json:"runtimeType,omitempty"`
	DockerPath     string                 `json:"dockerPath,omitempty"`
	DockerSum      string                 `json:"dockerSum,omitempty"`
	ContainerdPath string                 `json:"containerdPath,omitempty"`
	ContainerdSum  string                 `json:"containerdSum,omitempty"`
	NetworkType    string                 `json:"networkType,omitempty"`
	NetworkPath    string                 `json:"networkPath,omitempty"`
	NetworkSum     string                 `json:"networkSum,omitempty"`
	CalicoPath     string                 `json:"calicoPath,omitempty"`
	CalicoSum      string                 `json:"calicoSum,omitempty"`
	CiliumPath     string                 `json:"ciliumPath,omitempty"`
	CiliumSum      string                 `json:"ciliumSum,omitempty"`
	ContourPath    string                 `json:"contourPath,omitempty"`
	ContourSum     string                 `json:"contourSum,omitempty"`
	IstioPath      string                 `json:"istioPath,omitempty"`
	IstioSum       string                 `json:"istioSum,omitempty"`
	NvidiaPath     string                 `json:"nvidiaPath,omitempty"`
	NvidiaSum      string                 `json:"nvidiaSum,omitempty"`
	KataPath       string                 `json:"kataPath,omitempty"`
	KataSum        string                 `json:"kataSum,omitempty"`
	KruisePath     string                 `json:"kruisePath,omitempty"`
	KruiseSum      string                 `json:"kruiseSum,omitempty"`
	ClusterConf    *ClusterConfResource   `json:"clusterConf,omitempty"`
	ClusterScript  *ClusterScriptResource `json:"clusterScript,omitempty"`
}

type ClusterConfResource struct {
	KubeadmInitDir     string `json:"kubeadmInitDir"`
	KubeletServicePath string `json:"kubeletServicePath"`
	KubeadmConfPath    string `json:"kubeadmConfPath"`
	HaproxyStaticPath  string `json:"haproxyStaticPath"`
	StartupServicePath string `json:"startupServicePath"`
	StartupScriptPath  string `json:"startupScriptPath"`
}

type ClusterScriptResource struct {
	PreparePath        string `json:"preparePath"`
	PrepareCentosPath  string `json:"prepareCentosPath"`
	PrepareDebianPath  string `json:"prepareDebianPath"`
	PrepareUbuntuPath  string `json:"prepareUbuntuPath"`
	DiscoverPath       string `json:"discoverPath"`
	DiscoverNvidiaPath string `json:"discoverNvidiaPath"`
	SystemVersionPath  string `json:"systemVersionPath"`
}

type ClusterRemoteResource struct {
	BaseDir        string                       `json:"baseDir"`
	ConfDir        string                       `json:"confDir"`
	TplDir         string                       `json:"tplDir"`
	ScriptDir      string                       `json:"scriptDir"`
	PatchDir       string                       `json:"patchDir"`
	DistDir        string                       `json:"distDir"`
	TmpDir         string                       `json:"tmpDir"`
	ImagesPath     string                       `json:"imagesPath"`
	PausePath      string                       `json:"pausePath"`
	BinaryPath     string                       `json:"binaryPath"`
	KubeletPath    string                       `json:"kubeletPath"`
	OfflinePath    string                       `json:"offlinePath"`
	RuntimeType    string                       `json:"runtimeType"`
	DockerPath     string                       `json:"dockerPath"`
	ContainerdPath string                       `json:"containerdPath"`
	NetworkType    string                       `json:"networkType"`
	NetworkPath    string                       `json:"networkPath"`
	CalicoPath     string                       `json:"calicoPath"`
	CiliumPath     string                       `json:"ciliumPath"`
	ContourPath    string                       `json:"contourPath"`
	IstioPath      string                       `json:"istioPath"`
	NvidiaPath     string                       `json:"nvidiaPath"`
	KataPath       string                       `json:"kataPath"`
	KruisePath     string                       `json:"kruisePath"`
	ClusterConf    *ClusterRemoteConfResource   `json:"clusterConf"`
	ClusterScript  *ClusterRemoteScriptResource `json:"clusterScript"`
}

type ClusterRemoteConfResource struct {
	KubeletInitPath    string `json:"kubeletInitPath"`
	KubeadmInitPath    string `json:"kubeadmInitPath"`
	HaproxyStaticPath  string `json:"haproxyStaticPath"`
	StartupServicePath string `json:"startupServicePath"`
	StartupScriptPath  string `json:"startupScriptPath"`
}

type ClusterRemoteScriptResource struct {
	PreparePath        string `json:"preparePath"`
	PrepareCentosPath  string `json:"prepareCentosPath"`
	PrepareDebianPath  string `json:"prepareDebianPath"`
	PrepareUbuntuPath  string `json:"prepareUbuntuPath"`
	DiscoverPath       string `json:"discoverPath"`
	DiscoverNvidiaPath string `json:"discoverNvidiaPath"`
	SystemVersionPath  string `json:"systemVersionPath"`
}
