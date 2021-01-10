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
	"bufio"
	"bytes"
	"strings"
)

func ReadStringFromBuff(buff *bytes.Buffer) (result string, err error) {
	reader := bufio.NewReader(buff)
	result, err = reader.ReadString('\n')
	if nil != err {
		return "", err
	}
	return strings.TrimSpace(result), nil
}

func ReadLinesFromBuff(buff *bytes.Buffer) (lines []string) {
	scanner := bufio.NewScanner(buff)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

