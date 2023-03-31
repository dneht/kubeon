/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package viewsupport

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "support",
		Aliases: []string{"sp"},
		Short:   "Print support versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, version := range define.SupportVersionList() {
				fmt.Println(version)
			}
			return nil
		},
	}
	return cmd
}
