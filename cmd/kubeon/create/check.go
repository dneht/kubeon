/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package create

import (
	"github.com/dneht/kubeon/pkg/define"
	"k8s.io/klog/v2"
	"strings"
)

func checkSupport(flags *flagpole, clusterVersion string) bool {
	var isSupport bool
	isSupport = define.IsSupportVersion(clusterVersion)
	if !isSupport {
		klog.Errorf("[check] Input version[%s] not support", clusterVersion)
		return false
	}
	isSupport = define.IsSupportRuntime(flags.InputCRIMode)
	if !isSupport {
		klog.Errorf("[check] Input cri[%s] not support", flags.InputCRIMode)
		return false
	}
	isSupport = define.IsSupportNetwork(flags.InputCNIMode)
	if !isSupport {
		klog.Errorf("[check] Input cni[%s] not support", flags.InputCNIMode)
		return false
	}
	isSupport = define.IsSupportIngress(flags.InputICMode)
	if !isSupport {
		klog.Errorf("[check] Input ingress[%s] not support", flags.InputICMode)
		return false
	}
	isSupport = define.IsSupportProxyMode(flags.InputProxyMode)
	if !isSupport {
		klog.Errorf("[check] Input proxy mode[%s] not support", flags.InputProxyMode)
		return false
	}
	isSupport = define.IsSupportCiliumPolicyMode(flags.CiliumPolicyMode)
	if !isSupport {
		klog.Errorf("[check] Input cilium policy mode[%s] not support", flags.CiliumPolicyMode)
		return false
	}
	return true
}

func checkConfigs(config string) []string {
	config = strings.TrimSpace(config)
	if "" == config {
		return []string{}
	}
	return strings.Split(config, ",")
}
