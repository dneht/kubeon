/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package deleterouter

import (
	"github.com/dneht/kubeon/pkg/cloud"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "delete-router CLUSTER_NAME",
		Aliases: []string{"dr"},
		Short:   "Delete router for cloud vpc",
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterName := args[0]
			cluster.InitConfig(clusterName)
			_, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}
			current := cluster.Current()
			cloud.DeleteRouterNow(current.AllNodes)
			return nil
		},
	}
	return cmd
}
