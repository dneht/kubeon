/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Full 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cp

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(3),
		Use: "cp [flags] NODE_SELECTOR SRC_PATH DEST_PATH\n" +
			"Args:\n" +
			"  NODE_SELECTOR can be one of:\n" +
			"    @all 			all the control-plane and worker nodes \n" +
			"    @cp* 			all the control-plane nodes \n" +
			"    @cp1 			the bootstrap-control plane node \n" +
			"    @cpN 			the secondary master nodes \n" +
			"    @w* 			all the worker nodes\n" +
			"    @lb 			the external load balancer\n" +
			"    @etcd 			the external etcd\n" +
			"    @name=name 	the node hostname\n" +
			"    @ip=ip 		the node ip",
		Short: "Copy local files to node",
		Long:  "",
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
	return doCopy(cmd, args)
}

func doCopy(cmd *cobra.Command, args []string) error {
	nodeSelector := args[0]
	srcPath := args[1]
	destPath := args[2]
	nodes, err := cluster.SelectNodes(nodeSelector)
	if nil != err {
		return err
	}
	for _, node := range nodes {
		log.Infof("[%s] start copy file", node.Addr())
		err = node.CopyTo(srcPath, destPath)
		if nil != err {
			return err
		}
		log.Infof("[%s] copy file[%s] complete", node.Addr(), destPath)
	}
	return nil
}
