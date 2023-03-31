/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
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
