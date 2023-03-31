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

func KubectlCreateNamespace(name string) error {
	return execute.NewLocalCmd("kubectl", "create",
		"--kubeconfig="+cluster.Current().AdminConfigPath, "ns", name).RunWithEcho()
}
