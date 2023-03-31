/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"crypto/md5"
	"fmt"
)

func MD5(data []byte) string {
	md := md5.New()
	md.Write(data)
	return fmt.Sprintf("%x", md.Sum(nil))
}
