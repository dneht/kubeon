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

package add

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/spf13/cobra"
	"net"
)

type flagpole struct {
	define.DefaultList
	define.MasterList
	define.WorkerList
	DryRun      bool
	WithMirror  bool
	WithOffline bool
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "add [flags]\n",
		Short:   "Add a new node",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().BoolVar(
		&flags.DryRun, "dry-run",
		false, "dry run",
	)
	cmd.Flags().BoolVar(
		&flags.WithMirror, "with-mirror",
		true, "download use mirror, if in cn please keep true",
	)
	cmd.Flags().BoolVar(
		&flags.WithOffline, "with-offline",
		false, "install use offline system package",
	)
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
	newNodes, err := cluster.InitAddNodes(flags.DefaultList, flags.MasterList, flags.WorkerList, flags.DryRun)
	if nil != err {
		return err
	}
	if flags.DryRun {
		return nil
	}

	err = preInstall(newNodes, flags.WithMirror)
	if nil != err {
		return err
	}
	err = joinNodes(newNodes)
	if nil != err {
		return err
	}
	return nil
}

func preInstall(newNodes cluster.NodeList, mirror bool) (err error) {
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

func joinNodes(newNodes cluster.NodeList) (err error) {
	err = module.SetupAddKubeadm(newNodes)
	if nil != err {
		return err
	}
	newMasters := cluster.GetMasterFromList(newNodes)
	isNewMaster := len(newMasters) > 0
	if isNewMaster {
		err = module.InstallInner(define.CalicoNetwork)
		if nil != err {
			return err
		}
	}
	err = module.ChangeLoadBalance(isNewMaster, newNodes)
	if nil != err {
		return err
	}
	return cluster.AddCompleteCluster()
}
