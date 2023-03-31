/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package redo

import (
	"github.com/dneht/kubeon/cmd/kubeon/redo/redohosts"
	"github.com/dneht/kubeon/cmd/kubeon/redo/redoimages"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "redo NEXT_NAME",
		Aliases: []string{"R"},
		Short:   "Redo some resource",
	}
	cmd.AddCommand(redohosts.NewCommand())
	cmd.AddCommand(redoimages.NewCommand())
	return cmd
}
