/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package version

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Print kubeon version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("kubeon version: %s\n", define.AppVersion)
			return nil
		},
	}
	return cmd
}
