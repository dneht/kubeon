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
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"io/ioutil"
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

func CreateResource(mirror bool) (err error) {
	return release.ProcessDownload(current.LocalResource, current.Version.Full, current.RuntimeMode,
		mirror, current.IsBinary, current.IsOffline)
}

func CreateCompleteCluster() (err error) {
	err = UpdateHost()
	if nil != err {
		log.Errorf("update node host error: %s", err.Error())
		return err
	}

	current.Status = StatusRunning
	log.Infof("now cluster is running, api server is %s:%d", current.LbDomain, current.LbPort)
	err = runConfig.WriteConfig()
	if nil != err {
		log.Error("save cluster config failed: " + err.Error())
	}
	return nil
}

func UpgradeCompleteCluster() (err error) {
	err = UpdateHost()
	if nil != err {
		log.Errorf("update node host error: %s", err.Error())
		return err
	}

	current.Status = StatusRunning
	log.Infof("now cluster[%s] upgrade complete, version is %s", current.Name, current.Version.Full)
	err = runConfig.WriteConfig()
	if nil != err {
		log.Error("save cluster config failed: " + err.Error())
	}
	return nil
}

func AddCompleteCluster() (err error) {
	err = UpdateHost()
	if nil != err {
		log.Errorf("update node host error: %s", err.Error())
		return err
	}

	log.Infof("add nodes complete, version is %s", current.Version.Full)
	return nil
}

func bootBase64(path string) string {
	boot := BootstrapNode()
	result, err := boot.Command("cat", path, "|", "base64", "-w", "0").RunAndResult()
	if nil != err {
		log.Warnf("get remote[%s] config failed: %s, result: %s", boot.Addr(), err, result)
		return ""
	}
	return strings.TrimSpace(result)
}

func writeKubeConfig(caCert string) error {
	caData, err := base64.StdEncoding.DecodeString(caCert)
	if nil != err {
		return err
	}

	// create the directory to contain the KUBECONFIG file.
	// 0755 is taken from client-go's config handling logic: https://github.com/kubernetes/client-go/blob/5d107d4ebc00ee0ea606ad7e39fd6ce4b0d9bf9e/tools/clientcmd/loader.go#L412
	err = os.MkdirAll(filepath.Dir(current.AdminConfigPath), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create kubeconfig output directory")
	}
	return ioutil.WriteFile(current.AdminConfigPath, caData, 0600)
}
