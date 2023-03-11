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

package addon

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cloud"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"net"
)

type flagpole struct {
	define.DefaultList
	define.MasterList
	define.WorkerList
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "addon CLUSTER_NAME [flags]\n",
		Aliases: []string{"add"},
		Short:   "Add a node to the cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().UintVar(
		&flags.DefaultPort, "default-port",
		22, "All node default port",
	)
	cmd.Flags().StringVar(
		&flags.DefaultUser, "default-user",
		"root", "All node default username",
	)
	cmd.Flags().StringVar(
		&flags.DefaultPassword, "default-passwd",
		"", "All node default password",
	)
	cmd.Flags().StringVar(
		&flags.DefaultPkFile, "default-pk-file",
		"", "All node default private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringVar(
		&flags.DefaultPkPassword, "default-pk-password",
		"", "All node default private key file password",
	)
	cmd.Flags().IPSliceVarP(
		&flags.MasterIPs, "master", "m",
		[]net.IP{}, "Multi master ips",
	)
	cmd.Flags().UintSliceVar(
		&flags.MasterPorts, "master-port",
		[]uint{}, "Multi master port, default is 22",
	)
	cmd.Flags().StringSliceVarP(
		&flags.MasterNames, "master-name", "M",
		[]string{}, "Multi master name, if using will go to set hostname, otherwise use hostname",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterLabels, "master-labels",
		[]string{}, "Multi master labels",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterUsers, "master-user",
		[]string{}, "Multi master username",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterPasswords, "master-passwd",
		[]string{}, "Multi master password",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterPkFiles, "master-pk-file",
		[]string{}, "Multi master private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterPkPasswords, "master-pk-password",
		[]string{}, "Multi master private key file password",
	)
	cmd.Flags().UintVar(
		&flags.MasterDefaultPort, "default-master-port",
		22, "Multi master default port",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultUser, "default-master-user",
		"root", "Multi master default username",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultPassword, "default-master-passwd",
		"", "Multi master default password",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultPkFile, "default-master-pk-file",
		"", "Multi master default private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringVar(
		&flags.MasterDefaultPkPassword, "default-master-pk-password",
		"", "Multi master default private key file password",
	)
	cmd.Flags().IPSliceVarP(
		&flags.WorkerIPs, "worker", "w",
		[]net.IP{}, "Multi worker ips",
	)
	cmd.Flags().UintSliceVar(
		&flags.WorkerPorts, "worker-port",
		[]uint{}, "Multi worker port, default is 22",
	)
	cmd.Flags().StringSliceVarP(
		&flags.WorkerNames, "worker-name", "W",
		[]string{}, "Multi worker name, if using will go to set hostname, otherwise use hostname",
	)
	cmd.Flags().StringSliceVar(
		&flags.MasterLabels, "worker-labels",
		[]string{}, "Multi worker labels",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerUsers, "worker-user",
		[]string{}, "Multi worker username",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerPasswords, "worker-passwd",
		[]string{}, "Multi worker password",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerPkFiles, "worker-pk-file",
		[]string{}, "Multi master private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringSliceVar(
		&flags.WorkerPkPasswords, "worker-pk-password",
		[]string{}, "Multi worker private key file password",
	)
	cmd.Flags().UintVar(
		&flags.WorkerDefaultPort, "default-worker-port",
		22, "Multi worker default port",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultUser, "default-worker-user",
		"root", "Multi worker default username",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultPassword, "default-worker-passwd",
		"", "Multi worker default password",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultPkFile, "default-worker-pk-file",
		"", "Multi worker default private key file path, default is ~/.ssh/id_rsa",
	)
	cmd.Flags().StringVar(
		&flags.WorkerDefaultPkPassword, "default-worker-pk-password",
		"", "Multi worker default private key file password",
	)
	return cmd
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	_, err := cluster.InitExistCluster()
	if nil != err {
		return err
	}
	newNodes, err := cluster.InitAddNodes(flags.DefaultList, flags.MasterList, flags.WorkerList)
	if nil != err {
		return err
	}

	current := cluster.Current()
	err = preInstall(newNodes, current.Mirror)
	if nil != err {
		klog.Warningf("Prepare install failed, please check: %v", err)
		return nil
	}
	err = joinNodes(newNodes, current.IsOnCloud())
	if nil != err {
		klog.Errorf("Create nodes failed, reset nodes: %v", err)
		action.KubeadmResetList(newNodes, true, false)
	}
	return nil
}

func preInstall(newNodes cluster.NodeList, mirror string) (err error) {
	err = cluster.CreateResource(mirror)
	if nil != err {
		return err
	}

	err = module.PrepareInstall(newNodes, false)
	if nil != err {
		return err
	}
	return nil
}

func joinNodes(newNodes cluster.NodeList, onCloud bool) (err error) {
	err = module.SetupAddsKubeadm(newNodes)
	if nil != err {
		return err
	}
	if onCloud {
		cloud.ModifyRouterNow()
	}
	newMasters := cluster.GetMasterFromList(newNodes)
	isNewMaster := len(newMasters) > 0
	if isNewMaster {
		err = module.InstallNetwork(true)
		if nil != err {
			klog.Warningf("reinstall network failed %v", err)
		}
	}
	err = module.InstallExtend(false)
	if nil != err {
		klog.Warningf("reinstall extend failed %v", err)
	}
	err = module.ChangeLoadBalance(isNewMaster, newNodes)
	if nil != err {
		return err
	}
	module.LabelDevice()
	err = module.InstallIngress(false)
	if nil != err {
		klog.Warningf("reinstall ingress failed %v", err)
	}
	return cluster.AddCompleteCluster()
}
