/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/execute"
	"k8s.io/klog/v2"
)

const patchCorednsJson = "[{\"op\":\"add\",\"path\":\"/rules/-\",\"value\":{\"apiGroups\":[\"discovery.k8s.io\"],\"resources\":[\"endpointslices\"],\"verbs\":[\"list\",\"watch\"]}}]"

func KubectlPatchCorednsRole() error {
	output, err := KubectlGetResult(
		"clusterrole",
		"system:coredns",
		"-o=jsonpath='{.rules[-1].apiGroups}'",
	)
	if nil != err {
		return err
	}
	if output == "'[\"\"]'" {
		err = execute.NewLocalCmd("kubectl",
			"patch", "clusterrole", "system:coredns",
			"--type=json", "--patch="+patchCorednsJson,
			"--kubeconfig="+cluster.Current().AdminConfigPath,
		).RunWithEcho()
		if nil != err {
			klog.Warningf("Patch coredns failed: %v", err)
		}
	}
	return nil
}
