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
	"github.com/dneht/kubeon/pkg/onutil/log"
)

func DestroyCompleteCluster() (err error) {
	if nil == current {
		return nil
	}

	deleteResource()
	log.Infof("destroy cluster[%s] complete, version is %s", current.Name, current.Version.Full)
	return nil
}

func DelResetLocalHost(delNode *Node) {
	resetLocalHost(delNode)
}

func DelCompleteCluster(delNodes NodeList) (err error) {
	err = DeleteHost(delNodes)
	if nil != err {
		log.Errorf("delete node host error: %s", err.Error())
		return err
	}

	current.Status = StatusRunning
	log.Infof("delete nodes complete, version is %s", current.Version.Full)
	err = runConfig.WriteConfig()
	if nil != err {
		log.Error("del & save cluster config failed: " + err.Error())
	}
	return nil
}

func deleteResource() {
	DelConfig()
	onutil.RmDir(onutil.K8sDir())
}
