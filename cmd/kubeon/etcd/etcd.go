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

package etcd

import (
	"github.com/dneht/kubeon/cmd/kubeon/etcd/etcdbackup"
	"github.com/dneht/kubeon/cmd/kubeon/etcd/etcdcheck"
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.MinimumNArgs(2),
		Use:   "etcd CLUSTER_NAME ARGS...",
		Short: "Exec etcdctl",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cluster.InitExistCluster()
			if nil != err {
				return err
			}

			return action.Etcdctl(args[1:])
		},
	}
	cmd.AddCommand(etcdbackup.NewCommand())
	cmd.AddCommand(etcdcheck.NewCommand())
	return cmd
}
