/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/execute"
)

func KubectlLabelNodeRoleMaster(name string) error {
	return KubectlLabelRole(name, "node-role.kubernetes.io/master=")
}

func KubectlLabelNodeRoleWorker(name string) error {
	return KubectlLabelRole(name, "node-role.kubernetes.io/worker=")
}

func KubectlLabelRole(name, label string) error {
	labelArgs := []string{
		"label", "nodes", name, label,
		"--kubeconfig=" + cluster.Current().AdminConfigPath,
		"--overwrite",
	}
	return execute.NewLocalCmd("kubectl", labelArgs...).RunWithEcho()
}
