/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"net"
	"strings"
)

func IP2StringArr(ips []net.IP) []string {
	arr := make([]string, len(ips))
	for i, one := range ips {
		arr[i] = one.String()
	}
	return arr
}

func IsDuplicateInStringArr(arr []string) bool {
	hash := make(map[string]bool, len(arr)/2+1)
	for _, one := range arr {
		if one == "" {
			continue
		}
		_, ok := hash[one]
		if ok {
			return true
		} else {
			hash[one] = true
		}
	}
	return false
}

func SplitStringSpace(str string) []string {
	arr := strings.Split(str, "\n")
	if len(arr) == 1 {
		arr = strings.Split(str, "\r")
	}
	res := make([]string, 0, len(arr))
	for _, one := range arr {
		ts := strings.TrimSpace(one)
		if "" != ts {
			res = append(res, ts)
		}
	}
	return res
}
