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

package exec

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

type flagpole struct {
	Params     []string
	WithResult bool
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(2),
		Use: "exec NODE_SELECTOR COMMAND [flags]\n" +
			"Args:\n" +
			"  NODE_SELECTOR can be one of:\n" +
			"    cluter@all			all the control-plane and worker nodes \n" +
			"    cluter@cp*			all the control-plane nodes \n" +
			"    cluter@cp1			the bootstrap-control plane node \n" +
			"    cluter@cpN			the secondary master nodes \n" +
			"    cluter@w*			all the worker nodes\n" +
			"    cluter@name=name		the node hostname\n" +
			"    cluter@ip=ip		the node ip",
		Short: "Execute command on remote node",

		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().StringSliceVarP(
		&flags.Params, "arg", "p",
		[]string{},
		"command params",
	)
	cmd.Flags().BoolVarP(
		&flags.WithResult, "result", "R",
		false,
		"show result",
	)
	return cmd
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	clusterName, nodeSelector, err := onutil.NodeSelector(args[0])
	if nil != err {
		return err
	}
	cluster.InitConfig(clusterName)
	_, err = cluster.InitExistCluster()
	if nil != err {
		return err
	}
	return doExec(flags, nodeSelector, args[1])
}

func doExec(flags *flagpole, nodeSelector, command string) error {
	nodes, err := cluster.SelectNodes(nodeSelector)
	if nil != err {
		return err
	}
	for _, node := range nodes {
		exec := node.Command(command, flags.Params...)
		if flags.WithResult {
			err = exec.RunWithEcho()
			if nil != err {
				return err
			}
		} else {
			err = exec.Run()
			if nil != err {
				return err
			}
		}
		klog.V(1).Infof("[%s] exec command[%s] complete", node.Addr(), command)
	}
	return nil
}
