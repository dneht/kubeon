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
	"github.com/dneht/kubeon/pkg/onutil/log"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

// waitNewControlPlaneNodeReady waits for a new control plane node reaching the target state after init/join
func waitNewControlPlaneNodeReady(current *cluster.Cluster, node *cluster.Node, wait time.Duration) error {
	log.Infof("waiting for control-plane Pods to become ready (timeout %s)", wait)
	if pass := waitFor(current, node, wait,
		nodeIsReady,
		staticPodIsReady("kube-apiserver"),
		staticPodIsReady("kube-controller-manager"),
		staticPodIsReady("kube-scheduler"),
	); !pass {
		return errors.New("timeout: Node and control-plane did not reach target state")
	}
	fmt.Println()
	return nil
}

func waitForPodsRunning(current *cluster.Cluster, node *cluster.Node, wait time.Duration, label string, replicas int) error {
	if pass := waitFor(current, node, wait,
		podsAreRunning(label, replicas),
	); !pass {
		return errors.New("timeout: Node and control-plane did not reach target state")
	}
	fmt.Println()
	return nil
}

// waitForNodePort waits for a nodePort to become ready
func waitForNodePort(current *cluster.Cluster, node *cluster.Node, wait time.Duration, nodePort string) error {
	log.Infof("waiting for NodePort %q to become ready (timeout %s)", nodePort, wait)
	if pass := waitFor(current, node, wait,
		nodePortIsReady(nodePort),
	); !pass {
		return errors.New("timeout: NodePort not ready")
	}
	fmt.Println()
	return nil
}

// waitNewWorkerNodeReady waits for a new control plane node reaching the target state after join
func waitNewWorkerNodeReady(current *cluster.Cluster, node *cluster.Node, wait time.Duration) error {
	log.Infof("waiting for Node to become Ready (timeout %s)", wait)
	if pass := waitFor(current, node, wait,
		nodeIsReady,
	); !pass {
		return errors.New("timeout: Node did not reach target state")
	}
	fmt.Println()
	return nil
}

// waitControlPlaneUpgraded waits for a control plane node reaching the target state after upgrade
func waitControlPlaneUpgraded(current *cluster.Cluster, node *cluster.Node, upgradeVersion *define.StdVersion, wait time.Duration) error {
	version := upgradeVersion.Full

	log.Infof("waiting for control-plane Pods to restart with the new version (timeout %s)", wait)
	if pass := waitFor(current, node, wait,
		staticPodHasVersion("kube-apiserver", version),
		staticPodHasVersion("kube-controller-manager", version),
		staticPodHasVersion("kube-scheduler", version),
	); !pass {
		return errors.New("timeout: control-plane did not reach target state")
	}
	fmt.Println()
	return nil
}

// waitKubeletUpgraded waits for a node reaching the target state after upgrade
func waitKubeletUpgraded(current *cluster.Cluster, node *cluster.Node, upgradeVersion *define.StdVersion, wait time.Duration) error {
	version := upgradeVersion.Full

	log.Infof("waiting for node to restart with the new version (timeout %s)", wait)
	if pass := waitFor(current, node, wait,
		nodeHasKubernetesVersion(version),
	); !pass {
		return errors.New("timeout: node did not reach target state")
	}
	fmt.Println()
	return nil
}

// waitKubeletHasRBAC waits for the kubelet to have access to the expected config map
// please note that this is a temporary workaround for a problem we are observing on upgrades while
// executing node upgrades immediately after control-plane upgrade.
func waitKubeletHasRBAC(current *cluster.Cluster, node *cluster.Node, upgradeVersion *define.StdVersion, wait time.Duration) error {
	log.Infof("waiting for kubelet RBAC validation - workaround (timeout %s)", wait)
	if pass := waitFor(current, node, wait,
		kubeletHasRBAC(upgradeVersion.Major, upgradeVersion.Minor),
	); !pass {
		return errors.New("timeout: Node did not reach target state")
	}
	fmt.Println()
	return nil
}

// try defines a function that test a condition to be waited for
type try func(current *cluster.Cluster, node *cluster.Node) bool

// waitFor implements the waiter core logic that is responsible for testing all the given contitions
// until are satisfied or a timeout are reached
func waitFor(current *cluster.Cluster, node *cluster.Node, timeout time.Duration, conditions ...try) bool {
	// if timeout is 0 or no conditions are defined, exit fast
	if timeout == time.Duration(0) {
		fmt.Println("Timeout set 0, skipping wait")
		return true
	}

	// sets the timeout timer
	timer := time.NewTimer(timeout)

	// runs all the conditions in parallel
	pass := make(chan bool)
	for _, cond := range conditions {
		// run the condition in a go routine until it pass
		go func(check try) {
			// creates an arbitrary skew before starting a wait loop
			time.Sleep(time.Second + time.Duration(rand.Intn(500)+500)*time.Millisecond)

			for {
				if check(current, node) {
					pass <- true
					break
				}
				// add a little delay + jitter before retry
				time.Sleep(2*time.Second + time.Duration(rand.Intn(500)+500)*time.Millisecond)
			}
		}(cond)
	}

	// wait for all the conditions to pass or for a timeout
	passed := 0
	for {
		select {
		case <-pass:
			passed++
			if passed == len(conditions) {
				return true
			}
		case <-timer.C:
			return false
		}
	}
}
