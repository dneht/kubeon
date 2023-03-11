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

package cloud

import (
	"github.com/dneht/kubeon/pkg/cloud/alibaba"
	"github.com/dneht/kubeon/pkg/cloud/common"
	"github.com/dneht/kubeon/pkg/cloud/tencent"
	"github.com/dneht/kubeon/pkg/cluster"
	"k8s.io/klog/v2"
)

func ModifyRouterNow() {
	current := cluster.Current()
	if nil == current {
		klog.Warningf("[cloud] Please initialize cluster first, skip modify router")
		return
	}
	nodeInfoList, err := common.AllNodeCloudInfo()
	if nil != err {
		klog.Warningf("[cloud] Get k8s node info failed: %v", err)
		return
	}
	_ = ModifyRouter(current.CloudProvider, current.CloudConf, nodeInfoList)
}

func ModifyRouter(provider string, cloudConf *cluster.CloudConf, nodeInfoList map[string]*cluster.NodeCloudInfo) error {
	if nil == cloudConf || "" == provider {
		klog.Warningf("[cloud] Cloud provider is not set, skip modify router")
		return nil
	}
	endpoint, routerTableIds := cloudConf.Endpoint, cloudConf.RouterTableIds
	secretId, secretKey := common.GetAccessKeyFromEnv(provider)
	if "" == secretId || "" == secretKey || "" == endpoint {
		klog.Warningf("[cloud] Cloud access key is not set, skip modify router")
		return nil
	}

	switch provider {
	case "tencent":
		if err := tencent.ModifyRouter(secretId, secretKey, endpoint, routerTableIds, nodeInfoList); nil != err {
			klog.Warningf("[cloud] Cloud modify router on tencent failed: %v", err)
		}
	case "aliyun", "alibaba":
		if err := alibaba.ModifyRouter(secretId, secretKey, endpoint, routerTableIds, nodeInfoList); nil != err {
			klog.Warningf("[cloud] Cloud modify router on alibaba failed: %v", err)
		}
	}
	return nil
}

func DeleteRouterNow(delList cluster.NodeList) {
	if len(delList) == 0 {
		return
	}
	nodeInfoList := make(map[string]*cluster.NodeCloudInfo)
	for _, delOne := range delList {
		nodeInfoList[delOne.IPv4] = &cluster.NodeCloudInfo{
			Name: delOne.Hostname,
			IP:   delOne.IPv4,
		}
	}
	current := cluster.Current()
	if nil == current {
		klog.Warningf("[cloud] Please initialize cluster first, skip modify router")
		return
	}
	_ = DeleteRouter(current.CloudProvider, current.CloudConf, nodeInfoList)
}

func DeleteRouter(provider string, cloudConf *cluster.CloudConf, nodeInfoList map[string]*cluster.NodeCloudInfo) error {
	if nil == cloudConf || "" == provider {
		klog.Warningf("[cloud] Cloud provider is not set, skip delete router")
		return nil
	}
	endpoint, routerTableIds := cloudConf.Endpoint, cloudConf.RouterTableIds
	secretId, secretKey := common.GetAccessKeyFromEnv(provider)
	if "" == secretId || "" == secretKey || "" == endpoint {
		klog.Warningf("[cloud] Cloud access key is not set, skip delete router")
		return nil
	}

	switch provider {
	case "tencent":
		if err := tencent.DeleteRouter(secretId, secretKey, endpoint, routerTableIds, nodeInfoList); nil != err {
			klog.Warningf("[cloud] Cloud delete router on tencent failed: %v", err)
		}
	case "aliyun", "alibaba":
		if err := alibaba.DeleteRouter(secretId, secretKey, endpoint, routerTableIds, nodeInfoList); nil != err {
			klog.Warningf("[cloud] Cloud delete router on alibaba failed: %v", err)
		}
	}
	return nil
}
