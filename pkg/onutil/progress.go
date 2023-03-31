/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"k8s.io/klog/v2"
	"strconv"
)

const (
	OneKBByte = 1024
	OneMBByte = 1024 * 1024
)

func ShowProgress(total, done int64, prompt string) {
	if total <= 0 {
		return
	}
	totalLength, totalUnit := toSizeFormat(total)
	doneLength, doneUnit := toSizeFormat(done)
	doneProgress := strconv.FormatFloat(float64(done*100)/float64(total), 'f', 2, 64) + "%"
	klog.V(1).Infof("[%s] total: %s%s | done: %s%s | progress: %s", prompt,
		totalLength, totalUnit, doneLength, doneUnit, doneProgress)
}

func toSizeFormat(length int64) (string, string) {
	isMb := length/OneMBByte > 1
	if isMb {
		return strconv.FormatFloat(float64(length)/OneMBByte, 'f', 2, 64), "MB"
	} else {
		return strconv.FormatFloat(float64(length)/OneKBByte, 'f', 2, 64), "KB"
	}
}
