/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"github.com/pkg/errors"
	"strings"
)

func IpFromAddr(addr string) (string, error) {
	ipAndPort := strings.Split(addr, ":")
	if len(ipAndPort) != 2 {
		return "", errors.Errorf("addr format [%s] error", addr)
	}
	return ipAndPort[0], nil
}
