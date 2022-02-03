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

package create

import (
	"github.com/dneht/kubeon/pkg/action"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/module"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"net"
	"os"
	"time"
)

type flagpole struct {
	define.DefaultList
	define.MasterList
	define.WorkerList
	MirrorHost       string
	OnlyCreate       bool
	UseOffline       bool
	ClusterLBDomain  string
	ClusterDNSDomain string
	ClusterMaxPods   uint32
	ClusterPortRange string
	ExternalLBIP     string
	ExternalLBPort   int32
	InnerLBMode      string
	NodeInterface    []string
	NetworkSVCCIDR   string
	NetworkPodCIDR   string
	InputProxyMode   string
	IPVSScheduler    string
	InputCRIMode     string
	InputCNIMode     string
	CalicoMode       string
	CalicoMTU        string
	InputICMode      string
	WithNvidia       bool
	WithKata         bool
	InputCertSANs    []string
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(2),
		Use:     "create CLUSTER_NAME CLUSTER_VERSION [flags]\n",
		Aliases: []string{"c", "new"},
		Short:   "Create a new cluster",
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
	cmd.Flags().BoolVarP(
		&flags.OnlyCreate, "only-create", "C",
		false, "create config only",
	)
	cmd.Flags().BoolVar(
		&flags.UseOffline, "use-offline",
		false, "install use offline system package",
	)
	cmd.Flags().StringVar(
		&flags.ClusterLBDomain, "lb-domain",
		define.DefaultClusterAPIDomain,
		"A domain for apiserver load balance, not ip",
	)
	cmd.Flags().StringVar(
		&flags.ClusterDNSDomain, "dns-domain",
		define.DefaultClusterDNSDomain,
		"Cluster dns domain",
	)
	cmd.Flags().Uint32Var(
		&flags.ClusterMaxPods, "max-pods",
		define.DefaultClusterMaxPods,
		"Kubelet max pods",
	)
	cmd.Flags().StringVar(
		&flags.ClusterPortRange, "port-range",
		define.DefaultClusterPortRange,
		"Service node port range",
	)
	cmd.Flags().StringSliceVar(
		&flags.NodeInterface, "interface",
		[]string{},
		"If the node has multiple network, specify one",
	)
	cmd.Flags().StringVar(
		&flags.NetworkSVCCIDR, "svc-cidr",
		define.DefaultSVCSubnet,
		"Use alternative range of IP address for service VIPs",
	)
	cmd.Flags().StringVar(
		&flags.NetworkPodCIDR, "pod-cidr",
		define.DefaultPodSubnet,
		"Specify range of IP addresses for the pod network",
	)
	cmd.Flags().StringVar(
		&flags.ExternalLBIP, "lb-ip",
		"",
		"External load balancer ip",
	)
	cmd.Flags().Int32Var(
		&flags.ExternalLBPort, "lb-port",
		define.DefaultClusterAPIPort,
		"External load balancer port",
	)
	cmd.Flags().StringVar(
		&flags.InnerLBMode, "lb-mode",
		define.DefaultClusterLBMode,
		"Inner load balancer mode, only(haproxy or updater)",
	)
	cmd.Flags().StringVar(
		&flags.InputProxyMode, "proxy",
		define.DefaultProxyMode,
		"Kube proxy mode, only ipvs or iptables",
	)
	cmd.Flags().StringVar(
		&flags.IPVSScheduler, "ipvs-scheduler",
		define.DefaultIPVSScheduler,
		"IPVS scheduler, like: rr wrr",
	)
	cmd.Flags().StringVar(
		&flags.InputCRIMode, "cri",
		define.DefaultRuntimeMode,
		"Runtime interface, only docker or containerd",
	)
	cmd.Flags().StringVar(
		&flags.InputCNIMode, "cni",
		define.DefaultNetworkMode,
		"Network interface, only none or calico",
	)
	cmd.Flags().StringVar(
		&flags.CalicoMode, "calico-mode",
		define.DefaultCalicoMode,
		"Calico ip pool mode, only ipip or vxlan",
	)
	cmd.Flags().StringVar(
		&flags.CalicoMTU, "calico-mtu",
		define.DefaultCalicoMTU,
		"MTU for calico tunnel device",
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
	cmd.Flags().StringSliceVar(
		&flags.InputCertSANs, "cert-san",
		[]string{},
		"Kubernetes api server CertSANs",
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

func preRunE(flags *flagpole, cmd *cobra.Command, args []string) error {
	inputVersion, err := define.NewStdVersion(args[1])
	if nil != err {
		return err
	}
	version := inputVersion.Full
	if !checkSupport(flags, version) ||
		!flags.MasterList.CheckMatch() || !flags.WorkerList.CheckMatch() {
		os.Exit(1)
	}
	flags.WithNvidia = flags.WithNvidia && flags.InputCRIMode == define.ContainerdRuntime && inputVersion.IsSupportNvidia()
	flags.WithKata = flags.WithKata && inputVersion.IsSupportKata()
	return cluster.InitNewCluster(&cluster.Cluster{
		Version:       inputVersion,
		Mirror:        onutil.ConvMirror(flags.MirrorHost, define.MirrorImageRepo),
		IsOffline:     flags.UseOffline,
		LbPort:        flags.ExternalLBPort,
		LbMode:        flags.InnerLBMode,
		LbDomain:      flags.ClusterLBDomain,
		DnsDomain:     flags.ClusterDNSDomain,
		MaxPods:       flags.ClusterMaxPods,
		PortRange:     flags.ClusterPortRange,
		SvcCIDR:       flags.NetworkSVCCIDR,
		PodCIDR:       flags.NetworkPodCIDR,
		NodeInterface: flags.NodeInterface,
		ProxyMode:     flags.InputProxyMode,
		IPVSScheduler: flags.IPVSScheduler,
		RuntimeMode:   flags.InputCRIMode,
		NetworkMode:   flags.InputCNIMode,
		CalicoMode:    flags.CalicoMode,
		CalicoMTU:     flags.CalicoMTU,
		IngressMode:   flags.InputICMode,
		UseNvidia:     flags.WithNvidia,
		UseKata:       flags.WithKata,
		CertSANs:      flags.InputCertSANs,
		Status:        cluster.StatusCreating,
	}, flags.ExternalLBIP, flags.DefaultList, flags.MasterList, flags.WorkerList)
}

func runE(flags *flagpole, cmd *cobra.Command, args []string) (err error) {
	current := cluster.Current()
	if nil == current {
		return errors.New("cluster create error")
	}
	if flags.OnlyCreate {
		return nil
	}

	err = preInstall(current, current.Mirror)
	if nil != err {
		klog.Warningf("prepare input nodes failed, please check: %v", err)
		return nil
	}
	err = initCluster(current)
	if nil != err {
		klog.Errorf("Create cluster failed, reset nodes: %v", err)
		action.KubeadmResetList(cluster.CurrentNodes(), false)
	}
	return nil
}

func preInstall(current *cluster.Cluster, mirror string) (err error) {
	err = cluster.CreateResource(mirror)
	if nil != err {
		return err
	}

	err = module.PrepareInstall(cluster.CurrentNodes(), false)
	if nil != err {
		return err
	}
	return nil
}

func initCluster(current *cluster.Cluster) (err error) {
	bootNode := cluster.BootstrapNode()
	err = module.SetupBootKubeadm(bootNode)
	if nil != err {
		return err
	}
	err = module.InstallInner(define.HealthzReader)
	if nil != err {
		return err
	}
	err = module.InstallNetwork()
	if nil != err {
		return err
	}
	err = module.InstallExtend()
	if nil != err {
		return err
	}
	err = action.KubeadmInitWait(4 * time.Minute)
	if nil != err {
		return err
	}

	joinNodes := cluster.JoinsNode()
	err = module.SetupJoinsKubeadm(joinNodes)
	if nil != err {
		return err
	}
	err = module.InstallLoadBalance(current.Workers)
	if nil != err {
		return err
	}
	module.LabelDevice()
	return cluster.CreateCompleteCluster()
}
