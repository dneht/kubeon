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

package cluster

import (
	"github.com/dneht/kubeon/pkg/onutil"
	"k8s.io/klog/v2"
)

func DestroyCompleteCluster() (err error) {
	if nil == current {
		return nil
	}

	deleteResource()
	klog.V(1).Infof("Destroy cluster[%s] complete, version is %s", current.Name, current.Version.Full)
	return nil
}

func DelResetLocalHost(delNode *Node) {
	resetLocalHost(delNode)
}

func DelCompleteCluster(delNodes NodeList) (err error) {
	err = DeleteHost(delNodes)
	if nil != err {
		klog.Errorf("Delete node host error: %s", err.Error())
		return err
	}

	current.Status = StatusRunning
	klog.V(1).Infof("Delete nodes complete, version is %s", current.Version.Full)
	err = runConfig.WriteConfig()
	if nil != err {
		klog.Error("Delete & Save cluster config failed: " + err.Error())
	}
	return nil
}

func deleteResource() {
	DelConfig()
	onutil.RmDir(onutil.K8sDir())
}
