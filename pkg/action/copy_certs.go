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
	"k8s.io/klog/v2"
	"path/filepath"
)

const etcKubernetes = "/etc/kubernetes"

func CopyCertificates(current *cluster.Cluster) error {
	controlPlanes := current.ControlPlanes
	if len(controlPlanes) > 1 {
		for _, node := range controlPlanes[1:] {
			if err := copyCertToNode(current, node); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyCertToNode(current *cluster.Cluster, node *cluster.Node) error {
	// define the list of necessary cluster certificates
	fileNames := []string{
		"ca.crt", "ca.key",
		"front-proxy-ca.crt", "front-proxy-ca.key",
		"sa.pub", "sa.key",
	}
	if !current.IsExternalEtcd {
		fileNames = append(fileNames, "etcd/ca.crt", "etcd/ca.key")
	}

	if err := copyCertAndConfToNode(current, node, "pki", fileNames); err != nil {
		return err
	}
	return nil
}

func copyCAToNode(current *cluster.Cluster, node *cluster.Node) error {
	fileNames := []string{"ca.crt", "ca.key"}
	if err := copyCertAndConfToNode(current, node, "pki", fileNames); err != nil {
		return err
	}
	return nil
}

func copyConfigFilesToNode(current *cluster.Cluster, node *cluster.Node) error {
	fileNames := []string{
		"admin.conf",
		"controller-manager.conf",
		"scheduler.conf",
	}

	if err := copyCertAndConfToNode(current, node, "", fileNames); err != nil {
		return err
	}
	return nil
}

func copyCertAndConfToNode(current *cluster.Cluster, node *cluster.Node, basePath string, fileNames []string) error {
	klog.V(1).Infof("Importing cluster certificates from %s", cluster.BootstrapNode())

	for _, fileName := range fileNames {
		fmt.Printf("%s\n", fileName)
		filePath := filepath.Join(etcKubernetes, basePath, fileName)
		err := node.CopyTo(filePath, filePath)
		if nil != err {
			klog.Errorf("Copy cert|config[%s] to node[%s] failed", fileName, node.Addr())
			return err
		}
	}
	return nil
}
