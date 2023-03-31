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
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes", name,
		"node-role.kubernetes.io/master-",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes", name,
		"node-role.kubernetes.io/control-plane-",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
}

func KubectlRemoveAllMasterTaint() {
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes",
		"node-role.kubernetes.io/master-",
		"--all",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
	_ = execute.NewLocalCmd("kubectl",
		"taint", "nodes",
		"node-role.kubernetes.io/control-plane-",
		"--all",
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).Run()
}
