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

package kubeon

import (
	"github.com/dneht/kubeon/cmd/kubeon/add"
	"github.com/dneht/kubeon/cmd/kubeon/change"
	"github.com/dneht/kubeon/cmd/kubeon/cp"
	"github.com/dneht/kubeon/cmd/kubeon/create"
	"github.com/dneht/kubeon/cmd/kubeon/del"
	"github.com/dneht/kubeon/cmd/kubeon/destroy"
	"github.com/dneht/kubeon/cmd/kubeon/exec"
	"github.com/dneht/kubeon/cmd/kubeon/module"
	"github.com/dneht/kubeon/cmd/kubeon/upgrade"
	"github.com/dneht/kubeon/cmd/kubeon/version"
	"github.com/dneht/kubeon/cmd/kubeon/view"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"os"

	"github.com/spf13/cobra"
)

type flagpole struct {
	ClusterName string
	LogLevel    string
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "kubeon",
		Short: "kubeon is a simple way to create k8s clusters",
		Long: "   _          _                      \n" +
			"  | | ___   _| |__   ___  ___  _ __  \n" +
			"  | |/ / | | | '_ \\ / _ \\/ _ \\| '_ \\ \n" +
			"  |   <| |_| | |_) |  __/ (_) | | | |\n" +
			"  |_|\\_\\\\__,_|_.__/ \\___|\\___/|_| |_|\n\n" +
			"kubeon is still a work in progress. Test It, Break It, Send feedback!",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags, cmd, args)
		},
		SilenceUsage: true,
		Version:      define.AppVersion,
	}
	cmd.PersistentFlags().StringVarP(
		&flags.ClusterName,
		"name", "N",
		define.DefaultClusterName,
		"cluster name",
	)
	cmd.PersistentFlags().StringVar(
		&flags.LogLevel,
		"log-level",
		"",
		"log level, default(info)",
	)
	// add kubeon commands
	cmd.AddCommand(cp.NewCommand())
	cmd.AddCommand(exec.NewCommand())
	cmd.AddCommand(module.NewCommand())
	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(upgrade.NewCommand())
	cmd.AddCommand(destroy.NewCommand())
	cmd.AddCommand(add.NewCommand())
	cmd.AddCommand(del.NewCommand())
	cmd.AddCommand(change.NewCommand())
	cmd.AddCommand(view.NewCommand())
	cmd.AddCommand(version.NewCommand())
	return cmd
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) error {
	log.Init(flags.LogLevel)
	cluster.NewConfig(flags.ClusterName)
	return nil
}

// Run runs the root execute
func Run() error {
	return NewCommand().Execute()
}

// Main wraps Run and sets the log formatter
func Main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}
