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

package define

import (
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/pkg/errors"
	"strings"
)

type StdVersion struct {
	Full   string
	Number uint
	Major  uint
	Minor  uint
	Patch  uint
}

func NewStdVersion(version string) (*StdVersion, error) {
	if len(version) < 5 {
		return nil, errors.Errorf("input version [%s] error", version)
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
		return nil, errors.Errorf("input version [%s] error", version)
	}
	k8sVersion.Major = onutil.ParseUintOverZero(arr[0])
	k8sVersion.Minor = onutil.ParseUintOverZero(arr[1])
	k8sVersion.Patch = onutil.ParseUintOverZero(arr[2])
	k8sVersion.Number = k8sVersion.Major<<24 + k8sVersion.Minor<<16 + k8sVersion.Patch<<8
	return k8sVersion, nil
}

func (v *StdVersion) GreaterThen(in *StdVersion) bool {
	return v.Number > in.Number
}

func (v *StdVersion) LessThen(in *StdVersion) bool {
	return v.Number < in.Number
}
