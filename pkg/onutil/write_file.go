/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
	"errors"
	"os"
)

func WriteJsonFile(fileName, jsonData string) error {
	fileData, err := PrettyJson(jsonData)
	if nil != err {
		return err
	}
	return os.WriteFile(fileName, fileData, 0644)
}

func WriteFile(fileName string, fileData []byte) error {
	return os.WriteFile(fileName, fileData, 0644)
}

func ReadFile(filePath string) ([]byte, error) {
	if !PathExists(filePath) {
		return nil, errors.New("file " + filePath + " not exist")
	}
	return os.ReadFile(filePath)
}
