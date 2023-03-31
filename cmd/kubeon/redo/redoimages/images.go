/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package redoimages

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "images CLUSTER_NAME",
		Aliases: []string{"img", "image"},
		Short:   "Reimport container images",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}

			local := cluster.Current().IsRealLocal()
			for _, node := range cluster.CurrentNodes() {
				err = module.ImportImages(local, node)
				if nil != err {
					return err
				}
			}
			return nil
		},
	}
	return cmd
}
