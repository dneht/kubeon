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

func NodeSelector(name string) (string, string, error) {
	idx := strings.Index(name, "@")
	if idx <= 0 {
		return "", "", errors.New("selector format error")
	}
	if strings.Contains(name[idx:], "=") {
		return name[0:idx], name[idx+1:], nil
	} else {
		return name[0:idx], name[idx:], nil
	}
}
