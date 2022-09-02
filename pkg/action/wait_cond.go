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

package action

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"k8s.io/klog/v2"
	"strings"
	time "time"
)

// nodeIsReady implement a function that test when a node is ready
func nodeIsReady(current *cluster.Cluster, node *cluster.Node) bool {
	output, _ := KubectlGetQuiet(
		"nodes",
		// check for the selected node
		fmt.Sprintf("-l=kubernetes.io/hostname=%s", node.Hostname),
		// check for status.conditions type:Ready
		"-o=jsonpath='{.items..status.conditions[?(@.type==\"Ready\")].status}'",
	)
	if strings.Contains(output, "True") {
		klog.V(4).Infof("Node %s is ready\n", node.Hostname)
		return true
	} else {
		klog.V(6).Infof("Node %s jsonpath='{.items..status.conditions[?(@.type==\"Ready\")].status}' is not True\n", node.Hostname)
		return false
	}
}

// nodeHasKubernetesVersion implement a function that if a node is has the given Kubernetes version
func nodeHasKubernetesVersion(version string) func(current *cluster.Cluster, node *cluster.Node) bool {
	return func(current *cluster.Cluster, node *cluster.Node) bool {
		output, err := KubectlGetQuiet(
			"nodes",
			// check for the selected node
			fmt.Sprintf("-l=kubernetes.io/hostname=%s", node.Hostname),
			// check for the kubelet version
			"-o=jsonpath='{.items..status.nodeInfo.kubeletVersion}'",
		)
		if nil != err {
			return false
		}
		if strings.Contains(output, version) {
			klog.V(4).Infof("Node %s has Kubernetes version %s\n", node.Hostname, version)
			return true
		}
		return false
	}
}

// staticPodIsReady implement a function that test when a static pod is ready
func staticPodIsReady(pod string) func(current *cluster.Cluster, node *cluster.Node) bool {
	return func(current *cluster.Cluster, node *cluster.Node) bool {
		output, _ := KubectlGetQuiet(
			"pods",
			"-n=kube-system",
			// check for static pods existing on the selected node
			fmt.Sprintf("%s-%s", pod, node.Hostname),
			// check for status.conditions type:Ready
			"-o=jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}'",
		)
		if strings.Contains(output, "True") {
			klog.V(4).Infof("Pod %s-%s is ready\n", pod, node.Hostname)
			return true
		} else {
			klog.V(6).Infof("Pod %s-%s jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}' is not True\n", pod, node.Hostname)
			return false
		}
	}
}

func podsAreRunning(label string, replicas int) func(current *cluster.Cluster, node *cluster.Node) bool {
	return func(current *cluster.Cluster, node *cluster.Node) bool {
		output, err := KubectlGetQuiet(
			"pods",
			"-l", fmt.Sprintf("app=%s", label), "-o", "jsonpath='{.items[*].status.phase}'",
		)
		if nil != err {
			return false
		}
		statuses := strings.Split(strings.Trim(output, "'"), " ")

		// if pod number not yet converged, wait
		if len(statuses) != replicas {
			return false
		}

		// check for pods status
		running := true
		for j := 0; j < replicas; j++ {
			if statuses[j] != "Running" {
				running = false
			}
		}
		if running {
			klog.V(4).Infof("%d pods running!", replicas)
			return true
		}
		return false
	}
}

// nodePortIsReady implements a function that tests if a nodePort is ready
func nodePortIsReady(port string) func(current *cluster.Cluster, node *cluster.Node) bool {
	return func(current *cluster.Cluster, node *cluster.Node) bool {
		ip := node.IPv4
		result, err := node.Command(
			"curl", "-Is", fmt.Sprintf("http://%s:%s", ip, port),
		).RunAndResult()
		if err != nil {
			return false
		}
		if strings.Trim(result, "\n\r") == "HTTP/1.1 200 OK" {
			klog.V(4).Infof("Node port %s on node %s is ready...", port, node.Hostname)
			return true
		}
		return false
	}
}

// staticPodHasVersion implement a function that if a static pod is has the given Kubernetes version
func staticPodHasVersion(pod, version string) func(current *cluster.Cluster, node *cluster.Node) bool {
	return func(current *cluster.Cluster, node *cluster.Node) bool {
		output, err := KubectlGetQuiet(
			"pods",
			"-n=kube-system",
			// check for static pods existing on the selected node
			fmt.Sprintf("%s-%s", pod, node.Hostname),
			// check for the node image
			// NB. this assumes the Pod has only one container only
			// which is true for the control plane pods
			"-o=jsonpath='{.spec.containers[0].image}'",
		)
		if nil != err {
			return false
		}
		if strings.Contains(output, version) {
			klog.V(4).Infof("Pod %s-%s has Kubernetes version %s\n", pod, node.Hostname, version)
			return true
		}
		return false
	}
}

// kubeletHasRBAC is a test checking that kubelet has reliable access to kubelet-config-x.y and kube-proxy,
// where reliable = it have access for 5 seconds in a row.
//
// this test is a workaround meant to prevent errors like: configmaps "kube-proxy" is forbidden:
// User "system:node:kinder-upgrade-control-plane3" cannot get resource "configmaps" in API
// group "" in the namespace "kube-system": no relationship found between node "kinder-upgrade-control-plane3"
// and this object
//
// The real source of this errors during upgrades is still not clear, but it is probably related to
// the restarting of control-plane components after control-plane upgrade like e.g. the node authorizer
func kubeletHasRBAC(major, minor uint) func(current *cluster.Cluster, node *cluster.Node) bool {
	return func(current *cluster.Cluster, node *cluster.Node) bool {
		for i := 0; i < 5; i++ {
			output1, err := KubectlAuthGetResult(
				"--namespace=kube-system",
				fmt.Sprintf("configmaps/kubelet-config-%d.%d", major, minor),
			)
			if nil != err {
				return false
			}
			output2, err := KubectlAuthGetResult(
				"--namespace=kube-system",
				"configmaps/kube-proxy",
			)
			if nil != err {
				return false
			}
			if output1 == "yes" && output2 == "yes" {
				time.Sleep(1 * time.Second)
				continue
			}
			return false
		}

		klog.V(4).Infof("Kubelet has access to expected config maps")
		return true
	}
}
