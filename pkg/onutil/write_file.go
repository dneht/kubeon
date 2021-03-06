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
	"errors"
	"io/ioutil"
)

func WriteJsonFile(fileName, jsonData string) error {
	fileData, err := PrettyJson(jsonData)
	if nil != err {
		return err
	}
	return ioutil.WriteFile(fileName, fileData, 0644)
}

func WriteFile(fileName string, fileData []byte) error {
	return ioutil.WriteFile(fileName, fileData, 0644)
}

func ReadFile(filePath string) ([]byte, error) {
	if !PathExists(filePath) {
		return nil, errors.New("file " + filePath + " not exist")
	}
	return ioutil.ReadFile(filePath)
}
