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

package destroy

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "destroy CLUSTER_NAME [flags]\n",
		Aliases: []string{"stop"},
		Short:   "Destroy an exist cluster",
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

	err = resetCluster()
	if nil != err {
		log.Warnf("reset cluster failed, continue uninstall: %v", err)
	}
	err = doUninstall()
	if nil != err {
		log.Warnf("uninstall module failed, please check: %v", err)
	}
	return nil
}

func resetCluster() (err error) {
	action.KubeadmResetList(cluster.CurrentNodes())
	return nil
}

func doUninstall() (err error) {
	err = module.AllUninstall(cluster.CurrentNodes(), true)
	if nil != err {
		return err
	}

	return cluster.DestroyCompleteCluster()
}
