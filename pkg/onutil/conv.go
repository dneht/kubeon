/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"strings"
)

func ConvMirror(input, mirror, office string) string {
	if len(input) == 0 || strings.EqualFold(input, "no") || strings.EqualFold(input, "false") {
		return office
	} else if strings.EqualFold(input, "yes") || strings.EqualFold(input, "true") {
		return mirror
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
