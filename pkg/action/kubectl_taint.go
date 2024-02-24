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

func KubectlRemoveMasterTaint(name string) {
	masterLabel := cluster.Current().GetKubernetesMasterTaint()
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes", name, masterLabel+"-",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Quiet().RunWithEcho()
}

func KubectlRemoveAllMasterTaint() {
	masterLabel := cluster.Current().GetKubernetesMasterTaint()
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes", masterLabel+"-",
		"--all", "--kubeconfig="+cluster.Current().AdminConfigPath,
	).Quiet().RunWithEcho()
}
