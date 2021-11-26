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
