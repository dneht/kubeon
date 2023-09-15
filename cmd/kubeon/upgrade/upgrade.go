/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package upgrade

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

type flagpole struct {
	MirrorHost string
	SetOffline string
	WithNvidia bool
	WithKata   bool
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "upgrade CLUSTER_NAME CLUSTER_VERSION [flags]\n",
		Aliases: []string{"U", "up"},
		Short:   "Upgrade an exist cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cluster.InitConfig(args[0])
			return preRunE(flags, cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().StringVar(
		&flags.MirrorHost, "mirror",
		"yes", "download use mirror, if in cn please keep true",
	)
	cmd.Flags().StringVar(
		&flags.SetOffline, "set-offline",
		"", "modify upgrade offline mode",
	)
	cmd.Flags().BoolVar(
		&flags.WithNvidia, "with-nvidia",
		true,
		"Install nvidia",
	)
	cmd.Flags().BoolVar(
		&flags.WithKata, "with-kata",
		false,
		"Install kata with Kata-deploy",
	)
	return cmd
}

func preRunE(flags *flagpole, cmd *cobra.Command, args []string) error {
	inputVersion, err := define.NewStdVersion(args[1])
	if nil != err {
		return err
	}
	version := inputVersion.Full
	if !checkSupport(version) {
		os.Exit(1)
	}
	_, err = cluster.InitExistCluster()
	if nil != err {
		return err
	}
	current := cluster.Current()
	if "" == current.Mirror {
		current.Mirror = onutil.ConvMirror(flags.MirrorHost, define.MirrorImageRepo, define.DockerImageRepo)
	}
	if "" != flags.SetOffline {
		switch flags.SetOffline {
		case define.OnlineModule, "false", "no":
			current.IsOffline = false
			break
		case define.OfflineModule, "true", "yes":
			current.IsOffline = true
			break
		}
	}
	current.UseNvidia = (current.UseNvidia || flags.WithNvidia) && current.RuntimeMode == define.ContainerdRuntime
	current.UseKata = current.UseKata || flags.WithKata
	return cluster.InitUpgradeCluster(inputVersion)
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) (err error) {
	current := cluster.Current()
	if nil == current {
		return errors.New("cluster create error")
	}

	klog.V(1).Info("Ready to check & prepare host, please wait a moment...")
	err = preUpgrade(current, current.Mirror)
	if nil != err {
		return err
	}
	err = upgradeCluster(current)
	if nil != err {
		return err
	}
	return nil
}

func preUpgrade(current *cluster.Cluster, mirror string) (err error) {
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
	err = module.InstallNetwork(true)
	if nil != err {
		klog.Warningf("Reinstall network failed %v", err)
		return err
	}
	err = module.SetupUpgradeKubeadm()
	if nil != err {
		return err
	}
	for _, node := range cluster.CurrentNodes() {
		err = action.KubectlDrainNode(node.Hostname, current.Version)
		if nil != err {
			return err
		}
		err = action.KubeadmUpgrade(node, false)
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
	err = module.InstallDevice(true)
	if nil != err {
		klog.Warningf("Reinstall device failed %v", err)
	}
	err = module.UpgradeLoadBalance()
	if nil != err {
		return err
	}
	err = module.InstallIngress(true)
	if nil != err {
		klog.Warningf("Reinstall ingress failed %v", err)
	}
	err = module.InstallExtend(true)
	if nil != err {
		klog.Warningf("Reinstall extend failed %v", err)
	}
	return cluster.UpgradeCompleteCluster()
}
