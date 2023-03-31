/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package upgrade

import (
	"github.com/dneht/kubeon/pkg/define"
	"k8s.io/klog/v2"
)

func checkSupport(inputVersion string) bool {
	var isSupport bool
	isSupport = define.IsSupportVersion(inputVersion)
	if !isSupport {
		klog.Errorf("[check] Input version[%s] not support", inputVersion)
		return false
	}
	return isSupport
}
