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
	"github.com/dneht/kubeon/pkg/release"
	"github.com/spf13/cobra"
)

type flagpole struct {
	WithMirror  bool
	WithBinary  bool
	WithOffline bool
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Use: "download CLUSTER_VERSION\n" +
			"Args:\n" +
			"  CLUSTER_VERSION is you wanted kubernetes version",
		Aliases: []string{"d", "down"},
		Short:   "Download resource",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
	}
	cmd.Flags().BoolVar(
		&flags.WithMirror, "with-mirror",
		true, "download use mirror, if you is cn please keep true",
	)
	cmd.Flags().BoolVar(
		&flags.WithBinary, "with-binary",
		false, "download binary package",
	)
	cmd.Flags().BoolVarP(
		&flags.WithOffline, "with-offline", "O",
		false, "download offline system package",
	)
	return cmd
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	clusterVersion := args[0]
	runtimeMode := ""
	return release.ProcessDownload(release.InitResource(clusterVersion, runtimeMode, flags.WithBinary, flags.WithOffline),
		clusterVersion, runtimeMode, flags.WithMirror, flags.WithBinary, flags.WithOffline)
}
