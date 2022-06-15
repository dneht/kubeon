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
	"strings"
)

func ConvMirror(input, define string) string {
	if len(input) == 0 || strings.EqualFold(input, "no") || strings.EqualFold(input, "false") {
		return "no"
	} else if strings.EqualFold(input, "yes") || strings.EqualFold(input, "true") {
		return define
	} else {
		if strings.HasSuffix(input, "/") {
			return input[0 : len(input)-1]
		} else {
			return input
		}
	}
}

func ConvBool(res bool) string {
	if res {
		return "yes"
	} else {
		return "no"
	}
}
