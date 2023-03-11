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

package deleterouter

import (
	"github.com/dneht/kubeon/pkg/cloud"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "delete-router CLUSTER_NAME",
		Aliases: []string{"dr"},
		Short:   "Delete router for cloud vpc",
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterName := args[0]
			cluster.InitConfig(clusterName)
			_, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}
			current := cluster.Current()
			cloud.DeleteRouterNow(current.AllNodes)
			return nil
		},
	}
	return cmd
}
