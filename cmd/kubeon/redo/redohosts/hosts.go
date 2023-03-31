/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package redohosts

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "hosts CLUSTER_NAME",
		Aliases: []string{"host"},
		Short:   "Reconfigure /etc/hosts",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}

			err = cluster.InitHost()
			if nil != err {
				return err
			}
			return cluster.UpdateHost()
		},
	}
	return cmd
}
