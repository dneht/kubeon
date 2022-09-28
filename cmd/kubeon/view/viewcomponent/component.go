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
