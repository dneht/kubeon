/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Full 2.0 (the "License");
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
	"github.com/dneht/kubeon/pkg/onutil/log"
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
	log.Infof("[%s] total: %s%s | done: %s%s | progress: %s", prompt,
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
