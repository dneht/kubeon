/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
