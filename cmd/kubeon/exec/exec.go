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
	"sync"
)

type flagpole struct {
	Params      []string
	WithResult  bool
	UseParallel bool
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
		Short: "Execute command on the selected node",
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
	cmd.Flags().BoolVarP(
		&flags.UseParallel, "parallel", "P",
		false,
		"use parallel",
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
	if flags.UseParallel {
		var wait sync.WaitGroup
		wait.Add(len(nodes))
		for _, node := range nodes {
			go oneExec(flags, &wait, node, command)
		}
		wait.Wait()
	} else {
		for _, node := range nodes {
			oneExec(flags, nil, node, command)
		}
	}
	return nil
}

func oneExec(flags *flagpole, wait *sync.WaitGroup, node *cluster.Node, command string) {
	exec := node.Command(command, flags.Params...)
	if flags.WithResult {
		err := exec.RunWithEcho()
		if nil != err {
			klog.Errorf("[%s] exec command[%s] failed with: %v", node.Addr(), command, err)
		} else {
			klog.V(1).Infof("[%s] exec command[%s] complete", node.Addr(), command)
		}
	} else {
		err := exec.Run()
		if nil != err {
			klog.Errorf("[%s] exec command[%s] failed with: %v", node.Addr(), command, err)
		} else {
			klog.V(1).Infof("[%s] exec command[%s] complete", node.Addr(), command)

		}
	}
	if nil != wait {
		wait.Done()
	}
}
