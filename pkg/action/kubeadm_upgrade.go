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
	if current.Version.GreaterEqual(define.K8S_1_21_0) {
		return KubectlPatchCorednsRole()
	}
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
