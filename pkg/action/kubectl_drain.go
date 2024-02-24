/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"strconv"
)

func deleteLocalFlag(version *define.StdVersion) string {
	if version.GreaterEqual(define.K8S_1_20_0) {
		return "--delete-emptydir-data"
	} else {
		return "--delete-local-data"
	}
}

func KubectlDrainNode(name string, version *define.StdVersion, wait int) (err error) {
	return execute.NewLocalCmd("kubectl",
		"drain", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath,
		"--skip-wait-for-delete-timeout="+strconv.FormatInt(int64(wait), 10),
		"--ignore-daemonsets", deleteLocalFlag(version)).RunWithEcho()
}

func KubectlDrainNodeForce(name string, version *define.StdVersion, wait int) (err error) {
	return execute.NewLocalCmd("kubectl",
		"drain", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath,
		"--force",
		"--skip-wait-for-delete-timeout="+strconv.FormatInt(int64(wait), 10),
		"--ignore-daemonsets", deleteLocalFlag(version)).RunWithEcho()
}

func KubectlUncordonNode(name string) (err error) {
	return execute.NewLocalCmd("kubectl",
		"uncordon", name,
		"--kubeconfig="+cluster.Current().AdminConfigPath).RunWithEcho()
}
