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

package onutil

import (
	"io"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	maxRetry   = 3
	baseDomain = "dl.back.pub"
)

var baseDlUrl string

func init() {
	ips, _ := net.LookupIP(baseDomain)
	for _, ip := range ips {
		baseDlUrl = "http://" + ip.String() + "/on/"
		break
	}
	if len(baseDlUrl) < 15 {
		klog.Error("Network error, can not get remote server")
		os.Exit(1)
	}
}

func GetRemoteSum(version, moduleName string) string {
	return GetRemoteSumRetry(version, moduleName, 1)
}

func GetRemoteSumRetry(version, moduleName string, retry int) string {
	if retry > maxRetry {
		klog.Error("Network error, can not get remote version")
		os.Exit(1)
	}

	response, err := http.Get(baseDlUrl + version + "/" + moduleName + ".sum")
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
