/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/pkg/errors"
	"time"
)

func KubeadmJoinNode(nodes cluster.NodeList, usePatch bool, ignorePreflightErrors string, wait time.Duration) (err error) {
	current := cluster.Current()
	if usePatch && !current.Version.IsSupportPatch() {
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
	if usePatch && !current.Version.IsSupportPatch() {
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
	KubectlRemoveAllMasterTaint()
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
