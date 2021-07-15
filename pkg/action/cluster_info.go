/*
Copyright 2019 The Kubernetes Authors.

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

package action

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
)

// list of pods, pods images, etcd members
func ClusterInfo() error {
	current := cluster.Current()
	boot := cluster.BootstrapNode()
	// commands are executed on the bootstrap control-plane
	fmt.Println("====================cert info====================")
	if current.Version.GreaterThen(define.K8S_1_20_0) {
		if err := boot.Command(
			"kubeadm", "certs", "check-expiration",
		).RunWithEcho(); err != nil {
			return err
		}
	} else {
		if err := boot.Command(
			"kubeadm", "alpha", "certs", "check-expiration",
		).RunWithEcho(); err != nil {
			return err
		}
	}

	fmt.Println("====================node info====================")
	if err := boot.Command(
		"kubectl", "--kubeconfig=/etc/kubernetes/admin.conf", "get", "nodes", "-o=wide",
	).RunWithEcho(); err != nil {
		return err
	}

	fmt.Println("====================pod info====================")
	if err := boot.Command(
		"kubectl", "--kubeconfig=/etc/kubernetes/admin.conf", "get", "pods", "--all-namespaces", "-o=wide",
	).RunWithEcho(); err != nil {
		return err
	}

	fmt.Println("====================container info====================")
	if err := boot.Command(
		"kubectl", "--kubeconfig=/etc/kubernetes/admin.conf", "get", "pods", "--namespace=kube-system",
		"-o=jsonpath='{range .items[*]}{\"\\n\"}{.metadata.name}{\":\\t\"}{range .spec.containers[*]}{.image}{\", \"}{end}{end}'",
	).RunWithEcho(); err != nil {
		return err
	}
	fmt.Println()

	if !current.IsExternalEtcd {
		fmt.Println("====================etcd info====================")
		// Get the version of etcdctl from the etcd binary
		etcdctlVersion := EtcdVersion()
		fmt.Printf("using etcdctl version: %s\n", etcdctlVersion)
		err := etcdMemberListOutput(etcdctlVersion)
		if nil != err {
			return err
		}
	} else {
		fmt.Println("using external etcd")
	}
	return nil
}

func etcdMemberListOutput(version string) error {
	boot := cluster.BootstrapNode()
	memberArgs := buildEtcdctlArgs(boot)
	// Append version specific etcdctl certificate flags
	err := appendEtcdctlCertArgs(version, &memberArgs)
	if nil != err {
		return err
	}
	memberArgs = append(memberArgs, "member", "list", "--write-out=table")
	return boot.Command("kubectl", memberArgs...).RunWithEcho()
}
