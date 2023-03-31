/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package define

import "strings"

type ComponentVersion struct {
	Kubernetes string `json:"kubernetes"`
	Pause      string `json:"pause"`
	Etcd       string `json:"etcd"`
	Coredns    string `json:"coredns"`
	Crictl     string `json:"crictl,omitempty"`
	Runc       string `json:"runc,omitempty"`
	Containerd string `json:"containerd,omitempty"`
	Docker     string `json:"docker,omitempty"`
	Nvidia     string `json:"nvidia,omitempty"`
	Kata       string `json:"kata,omitempty"`
	Network    string `json:"cni,omitempty"`
	Calico     string `json:"calico,omitempty"`
	Cilium     string `json:"cilium,omitempty"`
	Hubble     string `json:"hubble,omitempty"`
	Contour    string `json:"contour,omitempty"`
	Istio      string `json:"istio,omitempty"`
	Haproxy    string `json:"haproxy,omitempty"`
	Kruise     string `json:"kruise,omitempty"`
	Offline    string `json:"offline"`
}

func (cv *ComponentVersion) RealNetwork() string {
	if "" == cv.Kubernetes {
		return ""
	}
	return cv.Kubernetes[:strings.LastIndex(cv.Kubernetes, ".")]
}

var SupportComponentFull = map[string]*ComponentVersion{}
