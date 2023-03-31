/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import "net"

var localAddresses = make(map[string]bool)

func init() {
	addresses, err := net.InterfaceAddrs()
	if nil == err {
		for _, addr := range addresses {
			ip, ok := addr.(*net.IPNet)
			if ok {
				localAddresses[ip.IP.String()] = true
			}
		}
	}
}

func IsLocalIPv4(ipv4 string) bool {
	_, ok := localAddresses[ipv4]
	return ok
}

func IsLocalIPv4InCluster(allIPs []string) bool {
	for _, ip := range allIPs {
		_, ok := localAddresses[ip]
		if ok {
			return true
		}
	}
	return false
}
