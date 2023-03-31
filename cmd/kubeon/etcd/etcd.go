/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package etcd

import (
	"github.com/dneht/kubeon/cmd/kubeon/etcd/etcdbackup"
	"github.com/dneht/kubeon/cmd/kubeon/etcd/etcdcheck"
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.MinimumNArgs(2),
		Use:   "etcd CLUSTER_NAME ARGS...",
		Short: "Exec etcdctl",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}

			return action.Etcdctl(args[1:])
		},
	}
	cmd.AddCommand(etcdbackup.NewCommand())
	cmd.AddCommand(etcdcheck.NewCommand())
	return cmd
}
