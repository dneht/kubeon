/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package viewconfig

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "config CLUSTER_NAME CONFIG_TYPE",
		Aliases: []string{"conf"},
		Short:   "Print select config details",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}

			moduleName := args[1]
			bytes, err := module.ShowInner(moduleName)
			if nil != err {
				return err
			}
			if nil == bytes {
				fmt.Printf("not found or not need this module %s", moduleName)
			} else {
				fmt.Println(string(bytes))
			}
			return nil
		},
	}
	return cmd
}
