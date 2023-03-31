/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package execute

import "strings"

func FileSum(path string) string {
	result, err := NewLocalCmd("sh", "-c", "if test -f "+path+"; then cksum "+path+"; else echo none; fi").RunAndResult()
	if nil != err || "none" == result {
		return ""
	} else {
		return strings.TrimSpace(strings.Split(result, " ")[0])
	}
}

func UnpackTar(filePath, unpackDir string) (err error) {
	return NewLocalCmd("tar", "xf", filePath, "-C", unpackDir).Run()
}

func UnpackTarGz(filePath, unpackDir string) (err error) {
	return NewLocalCmd("tar", "zxf", filePath, "-C", unpackDir).Run()
}

func UnpackTarXz(filePath, unpackDir string) (err error) {
	return NewLocalCmd("tar", "Jxf", filePath, "-C", unpackDir).Run()
}
