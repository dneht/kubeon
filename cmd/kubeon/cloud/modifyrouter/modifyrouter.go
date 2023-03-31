/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package modifyrouter

import (
	"github.com/dneht/kubeon/pkg/cloud"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/spf13/cobra"
)

type flagpole struct {
	CloudProvider       string
	CloudEndpoint       string
	CloudRouterTableIds []string
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "modify-router CLUSTER_NAME",
		Aliases: []string{"mr"},
		Short:   "Modify router for cloud vpc",
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterName := args[0]
			cluster.InitConfig(clusterName)
			config, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}
			current := cluster.Current()
			if "" != flags.CloudProvider {
				current.CloudProvider = flags.CloudProvider
			}
			if nil == current.CloudConf {
				current.CloudConf = &cluster.CloudConf{}
			}
			if "" != flags.CloudEndpoint {
				current.CloudConf.Endpoint = flags.CloudEndpoint
			}
			if len(flags.CloudRouterTableIds) > 0 {
				current.CloudConf.RouterTableIds = flags.CloudRouterTableIds
			}
			_ = config.WriteConfig()
			cloud.ModifyRouterNow()
			return nil
		},
	}
	cmd.Flags().StringVar(
		&flags.CloudProvider, "provider",
		"",
		"Cloud provider",
	)
	cmd.Flags().StringVar(
		&flags.CloudEndpoint, "endpoint",
		"",
		"Cloud endpoint",
	)
	cmd.Flags().StringSliceVar(
		&flags.CloudRouterTableIds, "router-table-id",
		[]string{},
		"Cloud router table id",
	)
	return cmd
}
