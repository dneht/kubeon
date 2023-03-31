/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package alibaba

import (
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/dneht/kubeon/pkg/cloud/common"
	"github.com/dneht/kubeon/pkg/cluster"
	"k8s.io/klog/v2"
	"strings"
)

func ModifyRouter(secretId, secretKey, endpoint string, routerTableIds []string, nodeInfos map[string]*cluster.NodeCloudInfo) error {
	remoteInfos, err := initRemoteRouterInfo(secretId, secretKey, endpoint, routerTableIds, nodeInfos)
	if nil != err {
		return err
	}
	if nil == remoteInfos || len(remoteInfos) == 0 {
		return nil
	}

	reNamePre, reDescPre := common.GetRouterNamePrefix(), common.GetRouterDescPrefix()
	for routerTableId, remoteInfoList := range remoteInfos {
		nodeInfoList := common.CopyNodeInfoList(nodeInfos)
		for _, remoteInfo := range remoteInfoList {
			instanceId := remoteInfo.InstanceId
			if nodeInfo, ok := nodeInfoList[instanceId]; ok {
				if remoteInfo.Desc == reDescPre+nodeInfo.Name && remoteInfo.CIDR == nodeInfo.CIDR {
					delete(nodeInfoList, instanceId)
				} else {
					nodeInfo.EntryId = remoteInfo.EntryId
				}
			}
		}
		delRouterEntries, addRouterEntries := make([]vpc.DeleteRouteEntriesRouteEntries, 0, len(nodeInfoList)), make([]vpc.CreateRouteEntriesRouteEntries, 0, len(nodeInfoList))
		for _, nodeInfo := range nodeInfoList {
			if "" == nodeInfo.InstanceId {
				continue
			}
			if "" != nodeInfo.EntryId {
				delRouterEntries = append(delRouterEntries, vpc.DeleteRouteEntriesRouteEntries{
					RouteTableId: routerTableId,
					RouteEntryId: nodeInfo.EntryId,
				})
			}
			addRouterEntries = append(addRouterEntries, vpc.CreateRouteEntriesRouteEntries{
				DstCidrBlock: nodeInfo.CIDR,
				RouteTableId: routerTableId,
				NextHop:      nodeInfo.InstanceId,
				NextHopType:  "Instance",
				Name:         reNamePre + nodeInfo.Name,
				Describption: reDescPre + nodeInfo.Name,
			})
		}
		if len(delRouterEntries) > 0 {
			request := vpc.CreateDeleteRouteEntriesRequest()
			request.Scheme = "https"
			request.RouteEntries = &delRouterEntries
			if _, err = vpcCli.DeleteRouteEntries(request); err != nil {
				return err
			}
		}
		if len(addRouterEntries) > 0 {
			request := vpc.CreateCreateRouteEntriesRequest()
			request.Scheme = "https"
			request.RouteEntries = &addRouterEntries
			if _, err = vpcCli.CreateRouteEntries(request); err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteRouter(secretId, secretKey, endpoint string, routerTableIds []string, nodeInfos map[string]*cluster.NodeCloudInfo) error {
	remoteInfos, err := initRemoteRouterInfo(secretId, secretKey, endpoint, routerTableIds, nodeInfos)
	if nil != err {
		return err
	}
	if nil == remoteInfos || len(remoteInfos) == 0 {
		return nil
	}

	reDescPre := common.GetRouterDescPrefix()
	for routerTableId, remoteInfoList := range remoteInfos {
		nodeInfoList := common.CopyNodeInfoList(nodeInfos)
		delRouterEntries := make([]vpc.DeleteRouteEntriesRouteEntries, 0, len(nodeInfoList))
		for _, remoteInfo := range remoteInfoList {
			if nodeInfo, ok := nodeInfoList[remoteInfo.InstanceId]; ok {
				if remoteInfo.Desc == reDescPre+nodeInfo.Name {
					delRouterEntries = append(delRouterEntries, vpc.DeleteRouteEntriesRouteEntries{
						RouteTableId: routerTableId,
						RouteEntryId: remoteInfo.EntryId,
					})
				}
			}
		}
		if len(delRouterEntries) > 0 {
			request := vpc.CreateDeleteRouteEntriesRequest()
			request.Scheme = "https"
			request.RouteEntries = &delRouterEntries
			if _, err = vpcCli.DeleteRouteEntries(request); err != nil {
				return err
			}
		}
	}
	return nil
}

func initRemoteRouterInfo(secretId, secretKey, endpoint string, routerTableIds []string,
	nodeInfoList map[string]*cluster.NodeCloudInfo) (map[string][]*cluster.NodeCloudInfo, error) {
	if len(nodeInfoList) == 0 {
		return nil, nil
	}
	if err := createClient(secretId, secretKey, endpoint); nil != err {
		return nil, err
	}

	privateIps := make([]string, len(nodeInfoList))
	for ip := range nodeInfoList {
		privateIps = append(privateIps, ip)
	}
	privateIpInfos, err := json.Marshal(privateIps)
	if nil != err {
		return nil, err
	}
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.InstanceNetworkType = "vpc"
	request.PrivateIpAddresses = string(privateIpInfos)
	response, err := ecsCli.DescribeInstances(request)
	if nil != err {
		return nil, err
	}
	for _, instance := range response.Instances.Instance {
		for _, ipAddr := range instance.VpcAttributes.PrivateIpAddress.IpAddress {
			if nodeInfo, ok := nodeInfoList[ipAddr]; ok {
				instanceId := instance.InstanceId
				nodeInfo.InstanceId = instanceId
				nodeInfoList[instanceId] = nodeInfo
				delete(nodeInfoList, ipAddr)
			}
		}
	}
	return getRemoteRouterInfo(routerTableIds)
}

func getRemoteRouterInfo(routerTableIds []string) (map[string][]*cluster.NodeCloudInfo, error) {
	current := cluster.Current()
	cloudInfos := make(map[string][]*cluster.NodeCloudInfo, len(routerTableIds))

	for _, routerTableId := range routerTableIds {
		cloudInfoList := make([]*cluster.NodeCloudInfo, 0, len(current.AllNodes))

		request := vpc.CreateDescribeRouteEntryListRequest()
		request.Scheme = "https"
		request.RouteTableId = routerTableId
		request.RouteEntryType = "Custom"
		request.NextHopType = "Instance"
		response, err := vpcCli.DescribeRouteEntryList(request)
		if err != nil {
			return nil, err
		}

		reDescPre := common.GetRouterDescPrefix()
		for _, routerEntry := range response.RouteEntrys.RouteEntry {
			instanceId := routerEntry.InstanceId
			if "" == instanceId && len(routerEntry.NextHops.NextHop) > 0 {
				instanceId = routerEntry.NextHops.NextHop[0].NextHopId
			}
			if "" != instanceId && strings.HasPrefix(routerEntry.Description, reDescPre) {
				klog.V(6).Infof("[cloud] Read router table[%s] entry: %s, %s, %v", routerTableId, routerEntry.RouteEntryId, instanceId, routerEntry.NextHops)
				cloudInfoList = append(cloudInfoList, &cluster.NodeCloudInfo{
					Desc:       routerEntry.Description,
					CIDR:       routerEntry.DestinationCidrBlock,
					InstanceId: routerEntry.InstanceId,
					EntryId:    routerEntry.RouteEntryId,
				})
			}
		}
		cloudInfos[routerTableId] = cloudInfoList
	}
	return cloudInfos, nil
}
