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

package copy

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"time"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.ExactArgs(3),
		Use: "copy NODE_SELECTOR SRC_PATH DEST_PATH [flags]\n" +
			"Args:\n" +
			"  NODE_SELECTOR can be one of:\n" +
			"    cluter@all 			all the control-plane and worker nodes \n" +
			"    cluter@cp* 			all the control-plane nodes \n" +
			"    cluter@cp1 			the bootstrap-control plane node \n" +
			"    cluter@cpN 			the secondary master nodes \n" +
			"    cluter@w* 			all the worker nodes\n" +
			"    cluter@name=name		the node hostname\n" +
			"    cluter@ip=ip		the node ip",
		Aliases: []string{"cp"},
		Short:   "Copy local files to node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(cmd, args)
		},
	}

	return cmd
}

func runE(cmd *cobra.Command, args []string) error {
	clusterName, nodeSelector, err := onutil.NodeSelector(args[0])
	if nil != err {
		return err
	}
	cluster.InitConfig(clusterName)
	_, err = cluster.InitExistCluster()
	if nil != err {
		return err
	}
	return doCopy(nodeSelector, args[1], args[2])
}

func doCopy(nodeSelector, srcPath, destPath string) error {
	nodes, err := cluster.SelectNodes(nodeSelector)
	if nil != err {
		return err
	}
	prog := mpb.New(
		mpb.WithWidth(90),
		mpb.WithRefreshRate(250*time.Millisecond),
	)
	for _, node := range nodes {
		oneCopy(prog, node, srcPath, destPath)
	}
	prog.Wait()
	return nil
}

func oneCopy(prog *mpb.Progress, node *cluster.Node, srcPath, destPath string) {
	bar := prog.New(0,
		mpb.BarStyle().Rbound("|"),
		mpb.PrependDecorators(
			decor.Name("copy to "+node.Hostname, decor.WC{W: 16, C: decor.DidentRight}),
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 0),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 0),
		),
	)
	go node.CopyToWithBar(srcPath, destPath, "", bar)
}
