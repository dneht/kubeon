/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package cloud

import (
	"github.com/dneht/kubeon/cmd/kubeon/cloud/deleterouter"
	"github.com/dneht/kubeon/cmd/kubeon/cloud/modifyrouter"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "cloud CMD_NAME",
		Aliases: []string{"cl"},
		Short:   "Cloud provider",
	}
	cmd.AddCommand(modifyrouter.NewCommand())
	cmd.AddCommand(deleterouter.NewCommand())
	return cmd
}
