/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package define

import (
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/pkg/errors"
	"strings"
)

type StdVersion struct {
	Full   string `json:"full"`
	Number uint   `json:"number"`
	Major  uint   `json:"major"`
	Minor  uint   `json:"minor"`
	Patch  uint   `json:"patch"`
}

func NewStdVersion(version string) (*StdVersion, error) {
	if len(version) < 5 {
		return nil, errors.Errorf("Input version [%s] error", version)
	}
	k8sVersion := &StdVersion{
		Full:   version,
		Number: 0,
	}
	var arr []string
	if strings.HasPrefix(version, "v") {
		arr = strings.Split(version[1:], ".")
	} else {
		arr = strings.Split(version, ".")
	}
	if len(arr) != 3 {
		return nil, errors.Errorf("Input version [%s] error", version)
	}
	k8sVersion.Major = onutil.ParseUintOverZero(arr[0])
	k8sVersion.Minor = onutil.ParseUintOverZero(arr[1])
	k8sVersion.Patch = onutil.ParseUintOverZero(arr[2])
	k8sVersion.Number = k8sVersion.Major<<24 + k8sVersion.Minor<<16 + k8sVersion.Patch<<8
	return k8sVersion, nil
}

func (v *StdVersion) IsSupportPatch() bool {
	return v.GreaterEqual(K8S_1_19_0)
}

func (v *StdVersion) IsSupportNvidia() bool {
	return v.GreaterEqual(K8S_1_22_0)
}

func (v *StdVersion) IsSupportKata() bool {
	return v.GreaterEqual(K8S_1_22_0)
}

func (v *StdVersion) GreaterThen(in *StdVersion) bool {
	return v.Number > in.Number
}

func (v *StdVersion) GreaterEqual(in *StdVersion) bool {
	return v.Number >= in.Number
}

func (v *StdVersion) LessThen(in *StdVersion) bool {
	return v.Number < in.Number
}

func (v *StdVersion) LessEqual(in *StdVersion) bool {
	return v.Number <= in.Number
}

func (v *StdVersion) String() string {
	return v.Full
}

type RngVersion struct {
	Start *StdVersion
	End   *StdVersion
}

func NewRngVersion(start string, end string) (*RngVersion, error) {
	if len(start) < 5 {
		return nil, errors.Errorf("Input version [%s] error", start)
	}
	startVersion, err := NewStdVersion(start)
	if nil != err {
		return nil, err
	}
	if len(end) < 5 {
		end = start
	}
	endVersion, err := NewStdVersion(end)
	if nil != err {
		return nil, err
	}
	if endVersion.Number < startVersion.Number {
		startVersion, endVersion = endVersion, startVersion
	}
	return &RngVersion{
		Start: startVersion,
		End:   endVersion,
	}, nil
}

func (v *RngVersion) Contain(other string) bool {
	input, err := NewStdVersion(other)
	if nil != err {
		return false
	}
	if input.Number <= v.End.Number && input.Number >= v.Start.Number {
		return true
	}
	return false
}

func (v *RngVersion) String() string {
	if v.Start.Number == v.End.Number {
		return v.Start.Full
	} else {
		return v.Start.Full + "-" + v.End.Full
	}
}
