/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package moduleuninstall

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "uninstall NODE_SELECTOR UNIT_NAME",
		Aliases: []string{"rm"},
		Short:   "Uninstall module",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd, args)
		},
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	clusterName, nodeSelector, err := onutil.NodeSelector(args[0])
	if nil != err {
		return err
	}
	cluster.InitConfig(clusterName)
	_, err = cluster.InitExistCluster()
	if nil != err {
		return err
	}

	moduleName := args[1]
	return module.UninstallSelect(moduleName, nodeSelector)
}
