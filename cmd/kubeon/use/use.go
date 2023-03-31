/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package use

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "use CLUSTER_NAME",
		Short: "Change default kubeconfig",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: runE,
	}
	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	config, err := cluster.InitExistCluster()
	if nil != err {
		return err
	}

	return config.ChangeConfig()
}
