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

package connect

import (
	"errors"
	"strconv"
)

var nodeConnect = make(map[string]*AuthConfig)

func SetAuthConfig(config *AuthConfig) (string, *AuthConfig, error) {
	if "" == config.Host {
		return "", nil, errors.New("remote host set error")
	}
	if config.Port == 0 {
		config.Port = 22
	}
	node := config.Host + ":" + strconv.FormatUint(config.Port, 10)
	nodeConnect[node] = config
	return node, config, nil
}

func GetAuthConfig(node string) *AuthConfig {
	return nodeConnect[node]
}
