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

package copy

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(3),
		Use: "copy NODE_SELECTOR SRC_PATH DEST_PATH [flags]\n" +
			"Args:\n" +
			"  NODE_SELECTOR can be one of:\n" +
			"    cluter@all 			all the control-plane and worker nodes \n" +
			"    cluter@cp* 			all the control-plane nodes \n" +
			"    cluter@cp1 			the bootstrap-control plane node \n" +
			"    cluter@cpN 			the secondary master nodes \n" +
			"    cluter@w* 			all the worker nodes\n" +
			"    cluter@name=name		the node hostname\n" +
			"    cluter@ip=ip		the node ip",
		Aliases: []string{"cp"},
		Short:   "Copy local files to node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd, args)
		},
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	clusterName, nodeSelector, err := onutil.NodeSelector(args[0])
	if nil != err {
		return err
	}
	cluster.InitConfig(clusterName)
	_, err = cluster.InitExistCluster()
	if nil != err {
		return err
	}
	return doCopy(nodeSelector, args[1], args[2])
}

func doCopy(nodeSelector, srcPath, destPath string) error {
	nodes, err := cluster.SelectNodes(nodeSelector)
	if nil != err {
		return err
	}
	for _, node := range nodes {
		klog.V(1).Infof("[%s] start copy file", node.Addr())
		err = node.CopyTo(srcPath, destPath)
		if nil != err {
			return err
		}
		klog.V(1).Infof("[%s] copy file[%s] complete", node.Addr(), destPath)
	}
	return nil
}
