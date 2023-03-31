/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"os"
	"os/user"
	"runtime"
)

func Home() string {
	current, err := user.Current()
	if nil == err {
		return current.HomeDir
	}

	// cross compile define
	if "windows" == runtime.GOOS {
		drive := os.Getenv("HOMEDRIVE")
		path := os.Getenv("HOMEPATH")
		home := drive + path
		if "" == drive || "" == path {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else {
		return os.Getenv("HOME")
	}
}

func BaseDir() string {
	return Home() + "/.kubeon"
}

func K8sDir() string {
	return Home() + "/.kube"
}
