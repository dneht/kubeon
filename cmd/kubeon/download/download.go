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
	InputICMode  string
	WithBinary   bool
	WithNvidia   bool
	WithKata     bool
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
		"yes", "default yes will use aliyun mirror, like: xxx.mirror.aliyuncs.com",
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
		&flags.InputICMode, "ic",
		define.DefaultIngressMode,
		"Ingress controller, only none or contour",
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

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	version, runtime, ingress := args[0], flags.InputCRIMode, flags.InputICMode
	resource := release.InitResource(version, runtime, flags.WithBinary, flags.UseOffline, flags.WithNvidia, flags.WithKata, ingress)
	return release.ProcessDownload(resource, version, runtime, onutil.ConvMirror(flags.MirrorHost, define.MirrorImageRepo),
		flags.ForceLocal, flags.WithBinary, flags.UseOffline, flags.WithNvidia, flags.WithKata, ingress)
}
