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
	"github.com/dneht/kubeon/pkg/onutil/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const maxRetry = 5
const baseDlUrl = "https://kubeon.oss-accelerate.aliyuncs.com/on/"

func GetRemoteSum(version, moduleName string) string {
	return GetRemoteSumRetry(version, moduleName, 1)
}

func GetRemoteSumRetry(version, moduleName string, retry int) string {
	if retry > maxRetry {
		log.Error("network error, can not get remote version")
		os.Exit(1)
	}

	response, err := http.Get(baseDlUrl + version + "/" + moduleName + ".sum")
	if nil != err {
		log.Debugf("get remote[%s -- %s] err is %s", version, moduleName, err)
		return GetRemoteSumRetry(version, moduleName, retry+1)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		log.Debugf("get remote[%s -- %s] err is %s", version, moduleName, err)
		return GetRemoteSumRetry(version, moduleName, retry+1)
	}
	return strings.TrimSpace(string(body))
}
