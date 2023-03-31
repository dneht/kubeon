/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"bytes"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/execute"
)

func KubectlApplyPath(path string) error {
	return execute.NewLocalCmd("kubectl", "apply",
		"--kubeconfig="+cluster.Current().AdminConfigPath, "-f", path).RunWithEcho()
}

func KubectlApplyData(content []byte) error {
	cmd := execute.NewLocalCmd("kubectl", "apply", "--kubeconfig="+cluster.Current().AdminConfigPath, "-f", "-")
	cmd.Stdin(bytes.NewBuffer(content))
	return cmd.RunWithEcho()
}

func KubectlDeleteData(content []byte) error {
	cmd := execute.NewLocalCmd("kubectl", "delete", "--kubeconfig="+cluster.Current().AdminConfigPath, "-f", "-")
	cmd.Stdin(bytes.NewBuffer(content))
	return cmd.RunWithEcho()
}
