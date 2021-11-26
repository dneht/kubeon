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

package upgrade

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"time"
)

type flagpole struct {
	WithMirror bool
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "upgrade CLUSTER_NAME CLUSTER_VERSION [flags]\n",
		Aliases: []string{"u", "up"},
		Short:   "Upgrade an exist cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return preRunE(flags, cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().BoolVar(
		&flags.WithMirror, "with-mirror",
		true, "download use mirror, if in cn please keep true",
	)
	return cmd
}

func preRunE(flags *flagpole, cmd *cobra.Command, args []string) error {
	inputVersion, err := define.NewStdVersion(args[1])
	if nil != err {
		return err
	}
	_, err = cluster.InitExistCluster()
	if nil != err {
		return err
	}
	return cluster.InitUpgradeCluster(inputVersion)
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) (err error) {
	current := cluster.Current()
	if nil == current {
		return errors.New("cluster create error")
	}

	err = preUpgrade(current, flags.WithMirror)
	if nil != err {
		return err
	}
	err = upgradeCluster(current)
	if nil != err {
		return err
	}
	return nil
}

func preUpgrade(current *cluster.Cluster, mirror bool) (err error) {
	err = cluster.CreateResource(mirror)
	if nil != err {
		return err
	}

	err = module.PrepareInstall(cluster.CurrentNodes(), true)
	if nil != err {
		return err
	}
	return nil
}

func upgradeCluster(current *cluster.Cluster) (err error) {
	//err = module.SetupUpgradeKubeadm()
	//if nil != err {
	//	return err
	//}
	for _, node := range cluster.CurrentNodes() {
		err = action.KubectlDrainNode(node.Hostname, current.Version)
		if nil != err {
			return err
		}
		if node.IsBootstrap() {
			err = action.KubeadmUpgradeApply(node, false, 4*time.Minute)
		} else {
			err = action.KubeadmUpgradeNode(node, false, 2*time.Minute)
		}
		if err != nil {
			return err
		}
		err = module.AfterUpgrade(node, node.IsBootstrap())
		if nil != err {
			return err
		}
		err = action.KubectlUncordonNode(node.Hostname)
		if nil != err {
			return err
		}
	}
	err = module.InstallNetwork()
	if nil != err {
		log.Warnf("reinstall network failed %v", err)
	}
	err = module.InstallIngress()
	if nil != err {
		log.Warnf("reinstall ingress failed %v", err)
	}
	err = module.UpgradeLoadBalance()
	if nil != err {
		return err
	}
	return cluster.UpgradeCompleteCluster()
}
