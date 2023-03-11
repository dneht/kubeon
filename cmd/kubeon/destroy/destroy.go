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
	"github.com/dneht/kubeon/pkg/cloud"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

type flagpole struct {
	ForceReset bool
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "destroy CLUSTER_NAME [flags]\n",
		Aliases: []string{"D", "rmr"},
		Short:   "Destroy an exist cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().BoolVarP(
		&flags.ForceReset, "force", "f",
		false, "if kubeadm reset hang, use force and reboot the machine after",
	)
	return cmd
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	_, err := cluster.InitExistCluster()
	if nil != err {
		return err
	}

	current := cluster.Current()
	err = resetCluster(flags.ForceReset)
	if nil != err {
		klog.Warningf("Reset cluster failed, continue uninstall: %v", err)
	}
	err = doUninstall(current)
	if nil != err {
		klog.Warningf("Uninstall module failed, please check: %v", err)
	}
	return nil
}

func resetCluster(force bool) (err error) {
	action.KubeadmResetList(cluster.CurrentNodes(), false, force)
	return nil
}

func doUninstall(current *cluster.Cluster) (err error) {
	err = module.AllUninstall(cluster.CurrentNodes(), true)
	if nil != err {
		return err
	}
	if current.IsOnCloud() {
		cloud.DeleteRouterNow(current.AllNodes)
	}
	return cluster.DestroyCompleteCluster()
}
