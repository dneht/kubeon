/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package alibaba

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

var vpcCli *vpc.Client
var ecsCli *ecs.Client

func createClient(secretId, secretKey, endpoint string) error {
	conf := sdk.NewConfig()
	cred := credentials.NewAccessKeyCredential(secretId, secretKey)
	if cli, err := ecs.NewClientWithOptions(endpoint, conf, cred); nil != err {
		return err
	} else {
		ecsCli = cli
	}
	if cli, err := vpc.NewClientWithOptions(endpoint, conf, cred); nil != err {
		return err
	} else {
		vpcCli = cli
	}
	return nil
}
