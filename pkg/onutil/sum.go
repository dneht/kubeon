/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"io"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strings"
)

const (
	maxRetry = 3
)

func GetRemoteSum(version, moduleName string) string {
	return GetRemoteSumRetry(version, moduleName, 0)
}

func GetRemoteSumRetry(version, moduleName string, retry int) string {
	if retry > maxRetry {
		klog.Error("Network error, can not get remote version")
		os.Exit(1)
	}

	baseDlUrl := "https://gitee.com/dneht/kubeon/raw/version/"
	if retry > 1 {
		baseDlUrl = "https://dl.back.pub/on/"
	}
	response, err := http.Get(baseDlUrl + moduleName + "/" + version + ".sum")
	if nil != err {
		klog.V(4).Infof("Get remote[%s -- %s] err is %s", version, moduleName, err)
		return GetRemoteSumRetry(version, moduleName, retry+1)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if nil != err {
		klog.V(4).Infof("Get remote[%s -- %s] err is %s", version, moduleName, err)
		return GetRemoteSumRetry(version, moduleName, retry+1)
	}
	return strings.TrimSpace(string(body))
}
