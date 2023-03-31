/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package tencent

import (
	"github.com/dneht/kubeon/pkg/cloud/common"
	"github.com/dneht/kubeon/pkg/cluster"
	util "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"k8s.io/klog/v2"
	"strconv"
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

	reDescPre := common.GetRouterDescPrefix()
	for routerTableId, remoteInfoList := range remoteInfos {
		nodeInfoList := common.CopyNodeInfoList(nodeInfos)
		for _, remoteInfo := range remoteInfoList {
			privateIp := remoteInfo.IP
			if nodeInfo, ok := nodeInfoList[privateIp]; ok {
				if remoteInfo.Desc == reDescPre+nodeInfo.Name && remoteInfo.CIDR == nodeInfo.CIDR {
					delete(nodeInfoList, privateIp)
				} else {
					nodeInfo.EntryId = remoteInfo.EntryId
					nodeInfo.RouterId = remoteInfo.RouterId
				}
			}
		}
		modRouterEntries, addRouterEntries := make([]*vpc.Route, 0, len(nodeInfoList)), make([]*vpc.Route, 0, len(nodeInfoList))
		for _, nodeInfo := range nodeInfoList {
			if nodeInfo.RouterId > 0 {
				modRouterEntries = append(modRouterEntries, &vpc.Route{
					DestinationCidrBlock: util.StringPtr(nodeInfo.CIDR),
					GatewayType:          util.StringPtr("NORMAL_CVM"),
					GatewayId:            util.StringPtr(nodeInfo.IP),
					RouteDescription:     util.StringPtr(reDescPre + nodeInfo.Name),
					RouteId:              util.Uint64Ptr(nodeInfo.RouterId),
				})
			} else {
				addRouterEntries = append(addRouterEntries, &vpc.Route{
					DestinationCidrBlock: util.StringPtr(nodeInfo.CIDR),
					GatewayType:          util.StringPtr("NORMAL_CVM"),
					GatewayId:            util.StringPtr(nodeInfo.IP),
					RouteDescription:     util.StringPtr(reDescPre + nodeInfo.Name),
				})
			}
		}
		if len(modRouterEntries) > 0 {
			request := vpc.NewReplaceRoutesRequest()
			request.RouteTableId = &routerTableId
			request.Routes = modRouterEntries
			if _, err = vpcCli.ReplaceRoutes(request); err != nil {
				klog.Errorf("mod: %v, %v, %v", routerTableId, len(modRouterEntries), err)
				return err
			}
		}
		if len(addRouterEntries) > 0 {
			request := vpc.NewCreateRoutesRequest()
			request.RouteTableId = &routerTableId
			request.Routes = addRouterEntries
			if _, err = vpcCli.CreateRoutes(request); err != nil {
				klog.Errorf("add: %v, %v, %v", routerTableId, len(addRouterEntries), err)
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
		delRouterEntries := make([]*vpc.Route, 0, len(nodeInfoList))
		for _, remoteInfo := range remoteInfoList {
			if nodeInfo, ok := nodeInfoList[remoteInfo.IP]; ok {
				if remoteInfo.Desc == reDescPre+nodeInfo.Name {
					delRouterEntries = append(delRouterEntries, &vpc.Route{
						RouteId: util.Uint64Ptr(remoteInfo.RouterId),
					})
				}
			}
		}
		if len(delRouterEntries) > 0 {
			request := vpc.NewDeleteRoutesRequest()
			request.RouteTableId = &routerTableId
			request.Routes = delRouterEntries
			if _, err = vpcCli.DeleteRoutes(request); err != nil {
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
	return getRemoteRouterInfo(routerTableIds)
}

func getRemoteRouterInfo(routerTableIds []string) (map[string][]*cluster.NodeCloudInfo, error) {
	current := cluster.Current()
	cloudInfos := make(map[string][]*cluster.NodeCloudInfo, len(routerTableIds))
	for _, routerTableId := range routerTableIds {
		cloudInfos[routerTableId] = make([]*cluster.NodeCloudInfo, 0, len(current.AllNodes))
	}

	request := vpc.NewDescribeRouteTablesRequest()
	request.RouteTableIds = util.StringPtrs(routerTableIds)
	request.Limit = util.StringPtr(strconv.FormatInt(int64(len(routerTableIds)), 10))
	response, err := vpcCli.DescribeRouteTables(request)
	if err != nil {
		return nil, err
	}

	reDescPre := common.GetRouterDescPrefix()
	for _, routerTable := range response.Response.RouteTableSet {
		routerTableId := *routerTable.RouteTableId
		cloudInfoList := cloudInfos[routerTableId]
		for _, routerEntry := range routerTable.RouteSet {
			privateIp := *routerEntry.GatewayId
			if "NORMAL_CVM" == *routerEntry.GatewayType && strings.HasPrefix(*routerEntry.RouteDescription, reDescPre) {
				klog.V(6).Infof("[cloud] Read router table[%s] entry: %s, %d, %s, %v", routerTableId, *routerEntry.RouteItemId, *routerEntry.RouteId, privateIp, *routerEntry.DestinationCidrBlock)
				cloudInfoList = append(cloudInfoList, &cluster.NodeCloudInfo{
					Desc:     *routerEntry.RouteDescription,
					IP:       privateIp,
					CIDR:     *routerEntry.DestinationCidrBlock,
					EntryId:  *routerEntry.RouteItemId,
					RouterId: *routerEntry.RouteId,
				})
			}
		}
		cloudInfos[routerTableId] = cloudInfoList
	}
	return cloudInfos, nil
}
