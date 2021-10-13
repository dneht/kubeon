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

package action

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/pkg/errors"
	"time"
)

func KubeadmJoinNode(nodes cluster.NodeList, usePatch bool, ignorePreflightErrors string, wait time.Duration) (err error) {
	current := cluster.Current()
	if usePatch && current.Version.LessThen(define.K8S_1_19_0) {
		return errors.New("--patches can't be used with kubeadm older than v1.19")
	}
	if len(nodes) > 0 {
		masterList := make(cluster.NodeList, 0, len(nodes))
		workerList := make(cluster.NodeList, 0, len(nodes))
		for _, node := range nodes {
			if node.IsControlPlane() {
				masterList = append(masterList, node)
			} else if node.IsWorker() {
				workerList = append(workerList, node)
			}
		}
		if len(masterList) > 0 {
			_, err = KubeadmJoinControlPlane(masterList, usePatch, ignorePreflightErrors, wait)
			if nil != err {
				return err
			}
		}
		if len(workerList) > 0 {
			_, err = KubeadmJoinWorker(workerList, ignorePreflightErrors, wait)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

func KubeadmJoinControlPlane(nodes cluster.NodeList, usePatch bool, ignorePreflightErrors string, wait time.Duration) (del *cluster.Node, err error) {
	current := cluster.Current()
	if usePatch && current.Version.LessThen(define.K8S_1_19_0) {
		return nil, errors.New("--patches can't be used with kubeadm older than v1.19")
	}
	for _, node := range nodes {
		err = kubeadmJoinControlPlane(node, usePatch, ignorePreflightErrors)
		if err != nil {
			return node, err
		}

		err = waitNewControlPlaneNodeReady(current, node, wait)
		if err != nil {
			return nil, err
		}
		err = joinControlPlaneAfterConfig(node)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func kubeadmJoinControlPlane(node *cluster.Node, usePatch bool, ignorePreflightErrors string) (err error) {
	joinArgs := []string{
		"join",
		fmt.Sprintf("--config=%s", node.GetResource().ClusterConf.KubeadmInitPath),
		fmt.Sprintf("--ignore-preflight-errors=%s", ignorePreflightErrors),
		fmt.Sprintf("--v=%d", log.Level()),
	}
	if usePatch {
		joinArgs = append(joinArgs, "--experimental-patches", node.GetResource().PatchDir)
	}
	return node.RunCmd("kubeadm", joinArgs...)
}

func joinControlPlaneAfterConfig(node *cluster.Node) error {
	return initAfterConfig(node)
}

func KubeadmJoinWorker(nodes cluster.NodeList, ignorePreflightErrors string, wait time.Duration) (del *cluster.Node, err error) {
	current := cluster.Current()
	for _, worker := range nodes {
		err = kubeadmJoinWorker(worker, ignorePreflightErrors)
		if err != nil {
			return worker, err
		}

		err = waitNewWorkerNodeReady(current, worker, wait)
		if err != nil {
			return nil, err
		}
		err = joinWorkerAfterConfig(worker)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func kubeadmJoinWorker(node *cluster.Node, ignorePreflightErrors string) (err error) {
	joinArgs := []string{
		"join",
		fmt.Sprintf("--config=%s", node.GetResource().ClusterConf.KubeadmInitPath),
		fmt.Sprintf("--ignore-preflight-errors=%s", ignorePreflightErrors),
		fmt.Sprintf("--v=%d", log.Level()),
	}
	return node.RunCmd("kubeadm", joinArgs...)
}

func joinWorkerAfterConfig(node *cluster.Node) error {
	return KubectlLabelNodeRoleWorker(node.Hostname)
}
