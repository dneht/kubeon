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

func KubectlDeleteNode(name string) error {
	return execute.NewLocalCmd("kubectl",
		"delete", "nodes", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath,
	).RunWithEcho()
}

func KubectlDeleteNodeForce(name string) error {
	return execute.NewLocalCmd("kubectl",
		"delete", "nodes", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath,
		"--force",
	).RunWithEcho()
}
