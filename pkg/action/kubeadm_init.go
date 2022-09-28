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

func KubeadmInitStart(boot *cluster.Node, uploadCerts, usePatch bool, ignorePreflightErrors string) (err error) {
	current := cluster.Current()
	if usePatch && !current.Version.IsSupportPatch() {
		return errors.New("--patches can't be used with kubeadm older than v1.19")
	}
	err = kubeadmInit(boot, uploadCerts, usePatch, ignorePreflightErrors)
	if nil != err {
		return err
	}
	return nil
}

func KubeadmInitWait(wait time.Duration) (err error) {
	boot := cluster.BootstrapNode()
	err = initAfterConfig(boot)
	if nil != err {
		return err
	}
	current := cluster.Current()
	if current.Version.GreaterEqual(define.K8S_1_21_0) {
		return KubectlPatchCorednsRole()
	}
	return waitNewControlPlaneNodeReady(cluster.Current(), boot, wait)
}

func kubeadmInit(boot *cluster.Node, uploadCerts, usePatch bool, ignorePreflightErrors string) (err error) {
	initArgs := []string{
		"init",
		fmt.Sprintf("--config=%s", boot.GetResource().ClusterConf.KubeadmInitPath),
		fmt.Sprintf("--ignore-preflight-errors=%s", ignorePreflightErrors),
		"--skip-token-print",
		"--skip-certificate-key-print",
		fmt.Sprintf("--v=%d", log.Level()),
	}
	if uploadCerts {
		initArgs = append(initArgs, "--upload-certs")
	}
	if usePatch {
		initArgs = append(initArgs, "--experimental-patches", boot.GetResource().PatchDir)
	}
	err = boot.RunCmd("kubeadm", initArgs...)
	if err != nil {
		return err
	}
	_, err = cluster.AfterBuildCluster()
	if err != nil {
		return err
	}
	return nil
}

func initAfterConfig(node *cluster.Node) error {
	KubectlRemoveMasterTaint(node.Hostname)
	return node.RunCmd("mkdir", "-p", node.Home+"/.kube",
		"&&", "\\cp", "/etc/kubernetes/admin.conf", node.Home+"/.kube/config",
		"&&", "chown", "$(id -u):$(id -g)", node.Home+"/.kube/config")
}

func KubeadmInitCert() (err error) {
	boot := cluster.BootstrapNode()
	return boot.RunCmd("kubeadm",
		"init", "phase", "certs", "all",
		"--config="+boot.GetResource().ClusterConf.KubeadmInitPath)
}

func KubeadmUploadCert(secretKey string) (err error) {
	boot := cluster.BootstrapNode()
	return boot.RunCmd("kubeadm",
		"init", "phase", "upload-certs",
		"--certificate-key="+secretKey,
		"--skip-certificate-key-print",
		"--upload-certs")
}
