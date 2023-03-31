/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import "strconv"

func ParseUintOverZero(in string) uint {
	get, err := strconv.ParseUint(in, 10, 64)
	if nil != err {
		return 0
	}
	return uint(get)
}
