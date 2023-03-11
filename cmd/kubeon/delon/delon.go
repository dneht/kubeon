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

package delon

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cloud"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "delon CLUSTER_NAME NODE_SELECTOR [flags]\n",
		Aliases: []string{"del", "rm"},
		Short:   "Remove a node from the cluster",
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
	delNodes, err := cluster.InitDelNodes(args[1])
	if nil != err {
		return err
	}

	current := cluster.Current()
	err = preRemove(delNodes)
	if nil != err {
		klog.Warningf("Prepare remove nodes failed, please check: %v", err)
		return nil
	}
	err = removeNodes(delNodes, current.IsOnCloud())
	if nil != err {
		klog.Errorf("Remove nodes failed, clean manually: %v", err)
	}
	return nil
}

func preRemove(delNodes cluster.NodeList) (err error) {
	for _, node := range delNodes {
		err = action.KubectlDrainNodeForce(node.Hostname, cluster.Current().Version)
		if nil != err {
			klog.Warningf("Drain node[%s] failed: %v", node.Addr(), err)
		}
	}
	return nil
}

func removeNodes(delNodes cluster.NodeList, onCloud bool) (err error) {
	action.KubeadmResetList(delNodes, true, false)
	err = module.AllUninstall(delNodes, false)
	if nil != err {
		return err
	}
	for _, node := range delNodes {
		cluster.DelResetLocalHost(node)
	}
	delMasters := cluster.GetMasterFromList(delNodes)
	isDelMaster := len(delMasters) > 0
	if isDelMaster {
		for _, node := range delMasters {
			err = action.EtcdMemberRemove(node.Hostname)
			if nil != err {
				klog.Warningf("Remove etcd member failed %v", err)
			}
		}
		err = module.InstallNetwork(true)
		if nil != err {
			klog.Warningf("Reinstall network failed %v", err)
		}
	}
	err = module.ChangeLoadBalance(isDelMaster, cluster.EmptyNodeList())
	if nil != err {
		return err
	}
	for _, node := range delNodes {
		err = action.KubectlDeleteNode(node.Hostname)
		if nil != err {
			klog.Warningf("Delete node[%s] failed: %v", node.Addr(), err)
		}
	}
	if onCloud {
		cloud.DeleteRouterNow(delNodes)
	}
	return cluster.DelCompleteCluster(delNodes)
}
