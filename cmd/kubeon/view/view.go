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
