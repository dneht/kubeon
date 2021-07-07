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
	"net"
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

func FillUintArr(one uint, size int) []uint {
	arr := make([]uint, size)
	for i := 0; i < size; i++ {
		arr[i] = one
	}
	return arr
}

func FillStringArr(one string, size int) []string {
	arr := make([]string, size)
	for i := 0; i < size; i++ {
		arr[i] = one
	}
	return arr
}
