/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package viewcomponent

import (
	"encoding/json"
	"fmt"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "component CLUSTER_VERSION",
		Aliases: []string{"cp"},
		Short:   "Print component install version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version := define.SupportComponentFull[args[0]]
			if nil == version {
				fmt.Println("input version " + args[0] + " is not support")
			} else {
				out, err := json.MarshalIndent(version, "", "    ")
				if nil != err {
					fmt.Println(err)
				} else {
					fmt.Println(string(out))
				}
			}
			return nil
		},
	}
	return cmd
}
