/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package view

import (
	"github.com/dneht/kubeon/cmd/kubeon/view/viewcomponent"
	"github.com/dneht/kubeon/cmd/kubeon/view/viewconfig"
	"github.com/dneht/kubeon/cmd/kubeon/view/viewsupport"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "view VIEW_NAME",
		Aliases: []string{"v"},
		Short:   "Print some info",
	}
	cmd.AddCommand(viewconfig.NewCommand())
	cmd.AddCommand(viewsupport.NewCommand())
	cmd.AddCommand(viewcomponent.NewCommand())
	return cmd
}
