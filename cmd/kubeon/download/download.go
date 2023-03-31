/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package download

import (
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/spf13/cobra"
)

type flagpole struct {
	MirrorHost   string
	ForceLocal   bool
	UseOffline   bool
	InputCRIMode string
	InputCNIMode string
	InputICMode  string
	WithBinary   bool
	WithNvidia   bool
	WithKata     bool
	WithKruise   bool
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Use: "download CLUSTER_VERSION\n" +
			"Args:\n" +
			"  CLUSTER_VERSION is you wanted kubernetes version",
		Aliases: []string{"down"},
		Short:   "Download install resources only",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().StringVar(
		&flags.MirrorHost, "mirror",
		"yes", "default yes will use cn mirror",
	)
	cmd.Flags().BoolVarP(
		&flags.ForceLocal, "force-local", "F",
		false, "install use local package",
	)
	cmd.Flags().BoolVar(
		&flags.UseOffline, "use-offline",
		false, "install use offline system package",
	)
	cmd.Flags().BoolVar(
		&flags.WithBinary, "with-binary",
		false, "download binary package",
	)
	cmd.Flags().StringVar(
		&flags.InputCRIMode, "cri",
		define.DefaultRuntimeMode,
		"Runtime interface, only docker or containerd",
	)
	cmd.Flags().StringVar(
		&flags.InputCNIMode, "cni",
		define.DefaultNetworkMode,
		"Network plugin, only none, calico or cilium",
	)
	cmd.Flags().StringVar(
		&flags.InputICMode, "ic",
		define.DefaultIngressMode,
		"Ingress controller, only none, contour or istio",
	)
	cmd.Flags().BoolVar(
		&flags.WithNvidia, "with-nvidia",
		true,
		"Install nvidia",
	)
	cmd.Flags().BoolVar(
		&flags.WithKata, "with-kata",
		false,
		"Install kata",
	)
	cmd.Flags().BoolVar(
		&flags.WithKruise, "with-kruise",
		false,
		"Install kruise",
	)
	return cmd
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	version, runtime, network, ingress := args[0], flags.InputCRIMode, flags.InputCNIMode, flags.InputICMode
	resource := release.InitResource(version, runtime, network, ingress,
		flags.WithBinary, flags.UseOffline, flags.WithNvidia, flags.WithKata, flags.WithKruise)
	return release.ProcessDownload(resource, version, runtime, network, ingress,
		onutil.ConvMirror(flags.MirrorHost, define.MirrorImageRepo, define.DockerImageRepo),
		flags.ForceLocal, flags.WithBinary, flags.UseOffline, flags.WithNvidia, flags.WithKata, flags.WithKruise)
}
