/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package moduleinner

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "inner CLUSTER_NAME UNIT_NAME",
		Aliases: []string{"in"},
		Short:   "Inner module",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd, args)
		},
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	cluster.InitConfig(args[0])
	_, err := cluster.InitExistCluster()
	if nil != err {
		return err
	}

	moduleName := args[1]
	return module.InstallInner(moduleName, false)
}
