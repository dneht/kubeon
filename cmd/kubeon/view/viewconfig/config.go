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

package viewconfig

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.ExactArgs(2),
		Use:   "config CLUSTER_NAME CONFIG_TYPE",
		Short: "Prints select config",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd, args)
		},
	}
	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	_, err := cluster.InitExistCluster()
	if nil != err {
		return err
	}

	moduleName := args[1]
	bytes, err :=module.ShowInner(moduleName)
	if nil != err {
		return err
	}
	if nil == bytes {
		fmt.Printf("not found or not need this module %s", moduleName)
	} else {
		fmt.Println(string(bytes))
	}
	return nil
}

