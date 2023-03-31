/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package module

import (
	"github.com/dneht/kubeon/cmd/kubeon/module/moduleinner"
	"github.com/dneht/kubeon/cmd/kubeon/module/moduleinstall"
	"github.com/dneht/kubeon/cmd/kubeon/module/moduleuninstall"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "module MODULE_NAME\n",
		Aliases: []string{"m"},
		Short:   "Install any module on any node",
		Long:    "",
		Example: "",
	}
	cmd.AddCommand(moduleinstall.NewCommand())
	cmd.AddCommand(moduleuninstall.NewCommand())
	cmd.AddCommand(moduleinner.NewCommand())
	return cmd
}
