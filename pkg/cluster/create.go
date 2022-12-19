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

package cluster

import (
	"encoding/base64"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"os"
	"path/filepath"
	"strings"
)

type CreateConfig struct {
	CACertBase64      string `json:"caCert"`
	EtcdKeyBase64     string `json:"etcdKey"`
	EtcdCertBase64    string `json:"etcdCert"`
	EtcdCABase64      string `json:"etcdCA"`
	EtcdEndpoints     string `json:"etcdEndpoints"`
	AdminConfigBase64 string `json:"adminConfig"`
}

func CreateResource(mirror string) error {
	return release.ProcessDownload(
		current.LocalResource, current.Version.Full,
		current.RuntimeMode, current.NetworkMode, current.IngressMode, mirror,
		current.IsRealLocal(), current.IsBinary, current.IsOffline,
		current.UseNvidia && current.HasNvidia, current.UseKata, current.UseKruise)
}

func CreateCompleteCluster() error {
	if err := UpdateHost(); nil != err {
		klog.Errorf("Update node host error: %s", err.Error())
		return err
	}

	current.Status = StatusRunning
	klog.V(1).Infof("Now cluster is running, api server is %s:%d", current.LbDomain, current.LbPort)
	if err := runConfig.WriteConfig(); nil != err {
		klog.Error("Create & Save cluster config failed: " + err.Error())
	}
	return nil
}

func UpgradeCompleteCluster() error {
	if err := UpdateHost(); nil != err {
		klog.Errorf("Update node host error: %s", err.Error())
		return err
	}

	current.Status = StatusRunning
	klog.V(1).Infof("Now cluster[%s] upgrade complete, version is %s", current.Name, current.Version.Full)
	if err := runConfig.WriteConfig(); nil != err {
		klog.Error("Upgrade & save cluster config failed: " + err.Error())
	}
	return nil
}

func AddCompleteCluster() error {
	if err := UpdateHost(); nil != err {
		klog.Errorf("Update node host error: %s", err.Error())
		return err
	}

	current.Status = StatusRunning
	klog.V(1).Infof("Add nodes complete, version is %s", current.Version.Full)
	if err := runConfig.WriteConfig(); nil != err {
		klog.Error("Add & Save cluster config failed: " + err.Error())
	}
	return nil
}

func bootBase64(path string) string {
	boot := BootstrapNode()
	result, err := boot.Command("cat", path, "|", "base64", "-w", "0").RunAndResult()
	if nil != err {
		klog.Warningf("Get remote[%s] config failed: %s, result: %s", boot.Addr(), err, result)
		return ""
	}
	return strings.TrimSpace(result)
}

func loadCreateConfig() string {
	adminConfigBase64 := bootBase64(define.KubernetesEtcPath + "/admin.conf")
	current.CreateConfig = &CreateConfig{
		CACertBase64:      bootBase64(define.KubernetesPkiPath + "/ca.crt"),
		EtcdKeyBase64:     bootBase64(define.KubernetesPkiPath + "/apiserver-etcd-client.key"),
		EtcdCertBase64:    bootBase64(define.KubernetesPkiPath + "/apiserver-etcd-client.crt"),
		EtcdCABase64:      bootBase64(define.KubernetesPkiEtcdPath + "/ca.crt"),
		EtcdEndpoints:     etcdEndpoints(),
		AdminConfigBase64: adminConfigBase64,
	}
	return adminConfigBase64
}

func writeKubeConfig(base64Data string) error {
	decodeData, err := base64.StdEncoding.DecodeString(base64Data)
	if nil != err {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(current.AdminConfigPath), 0755); err != nil {
		return errors.Wrap(err, "failed to create kubeconfig output directory")
	}
	return os.WriteFile(current.AdminConfigPath, decodeData, 0600)
}
