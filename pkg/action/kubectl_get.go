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

func KubectlGetQuiet(args ...string) (string, error) {
	return kubectlGet(args).Quiet().RunAndResult()
}

func KubectlGetResult(args ...string) (string, error) {
	return kubectlGet(args).RunAndResult()
}

func KubectlGetOutput(args ...string) error {
	return kubectlGet(args).RunWithEcho()
}

func kubectlGet(args []string) *execute.LocalCmd {
	getArgs := []string{
		"get",
		"--kubeconfig=" + cluster.Current().AdminConfigPath,
	}
	getArgs = append(getArgs, args...)
	return execute.NewLocalCmd("kubectl", getArgs...)
}

func KubectlAuthGetResult(args ...string) (string, error) {
	getArgs := []string{
		"auth",
		"can-i",
		"get",
		"--kubeconfig=" + cluster.Current().AdminConfigPath,
	}
	getArgs = append(getArgs, args...)
	return execute.NewLocalCmd("kubectl", getArgs...).RunAndResult()
}
