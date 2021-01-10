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

package onutil

import (
	"os"
	"os/user"
	"runtime"
)

func Home() string {
	current, err := user.Current()
	if nil == err {
		return current.HomeDir
	}

	// cross compile define
	if "windows" == runtime.GOOS {
		drive := os.Getenv("HOMEDRIVE")
		path := os.Getenv("HOMEPATH")
		home := drive + path
		if "" == drive || "" == path {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else {
		return os.Getenv("HOME")
	}
}

func BaseDir() string {
	return Home() + "/.kubeon"
}

func K8sDir() string {
	return Home() + "/.kube"
}
