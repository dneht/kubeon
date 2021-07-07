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

package del

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "del [flags] NODE_SELECTOR\n",
		Short:   "Del an exit node",
		Long:    "",
		Example: "",
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
	delNodes, err := cluster.InitDelNodes(args[0])
	if nil != err {
		return err
	}

	err = preRemove(delNodes)
	if nil != err {
		return err
	}
	err = removeNodes(delNodes)
	if nil != err {
		return err
	}
	return nil
}

func preRemove(delNodes cluster.NodeList) (err error) {
	for _, node := range delNodes {
		err = action.KubectlDrainNodeForce(node.Hostname)
		if nil != err {
			log.Warnf("drain node[%s] error: %s", node.Addr(), err)
		}
		err = action.KubectlDeleteNode(node.Hostname)
		if nil != err {
			log.Warnf("delete node[%s] error: %s", node.Addr(), err)
		}
	}
	action.KubeadmResetForce(delNodes)
	return nil
}

func removeNodes(delNodes cluster.NodeList) (err error) {
	err = module.AllUninstall(delNodes, false)
	if nil != err {
		return err
	}
	delMasters := cluster.GetMasterFromList(delNodes)
	isDelMaster := len(delMasters) > 0
	if isDelMaster {
		for _, node := range delMasters {
			err = action.EtcdMemberRemove(node.Hostname)
			if nil != err {
				return err
			}
		}
		err = module.InstallInner(define.CalicoNetwork)
		if nil != err {
			return err
		}
	}
	err = module.ChangeLoadBalance(isDelMaster, cluster.EmptyNodeList())
	if nil != err {
		return err
	}
	return cluster.DelCompleteCluster(delNodes)
}
