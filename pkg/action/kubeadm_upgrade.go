/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
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

func KubeadmUpgrade(node *cluster.Node, usePatch bool) (err error) {
	if node.IsBootstrap() {
		err = KubeadmUpgradeApply(node, usePatch, define.KubeadmIgnorePreflightErrors, 4*time.Minute)
	} else {
		err = KubeadmUpgradeNode(node, usePatch, define.KubeadmIgnorePreflightErrors, 2*time.Minute)
	}
	return err
}

func KubeadmUpgradeApply(boot *cluster.Node, usePatch bool, ignorePreflightErrors string, wait time.Duration) (err error) {
	current := cluster.Current()
	if usePatch && !current.Version.IsSupportPatch() {
		return errors.New("--patches can't be used with kubeadm older than v1.19")
	}
	applyArgs := []string{
		"upgrade", "apply", current.Version.Full,
		"--force",
		"--certificate-renewal=false",
		fmt.Sprintf("--config=%s", boot.GetResource().ClusterConf.KubeadmInitPath),
		fmt.Sprintf("--ignore-preflight-errors=%s", ignorePreflightErrors),
		fmt.Sprintf("--v=%d", log.Level()),
	}
	if usePatch {
		applyArgs = append(applyArgs, "--experimental-patches", define.AppPatchDir)
	}
	err = boot.RunCmd("kubeadm", applyArgs...)
	if err != nil {
		return err
	}
	_, err = cluster.AfterBuildCluster()
	if err != nil {
		return err
	}
	_ = KubectlPatchCorednsRole()
	return waitControlPlaneUpgraded(current, boot, current.Version, wait)
}

func KubeadmUpgradeNode(node *cluster.Node, usePatch bool, ignorePreflightErrors string, wait time.Duration) (err error) {
	current := cluster.Current()
	// waitKubeletHasRBAC waits for the kubelet to have access to the expected config map
	// please note that this is a temporary workaround for a problem we are observing on upgrades while
	// executing node upgrades immediately after control-plane upgrade.
	if err = waitKubeletHasRBAC(current, node, current.Version, wait); err != nil {
		return err
	}

	// kubeadm upgrade node
	nodeArgs := []string{
		"upgrade", "node",
		"--certificate-renewal=false",
		fmt.Sprintf("--ignore-preflight-errors=%s", ignorePreflightErrors),
		fmt.Sprintf("--v=%d", log.Level()),
	}
	if usePatch {
		nodeArgs = append(nodeArgs, "--experimental-patches", define.AppPatchDir)
	}
	err = node.RunCmd("kubeadm", nodeArgs...)
	if err != nil {
		return err
	}
	if node.IsControlPlane() {
		err = waitControlPlaneUpgraded(current, node, current.Version, wait)
		if err != nil {
			return err
		}
	}
	return nil
}
