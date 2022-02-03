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
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func PathIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func PathIsFile(path string) bool {
	return !PathIsDir(path)
}

func LsDir(path string) []os.FileInfo {
	list := make([]os.FileInfo, 0)
	if PathIsDir(path) {
		dir, err := ioutil.ReadDir(path)
		if nil == err {
			for _, one := range dir {
				if !one.IsDir() {
					list = append(list, one)
				}
			}
		} else {
			klog.Warningf("Read path[%s] error")
		}
	}
	return list
}

func MkDir(path string) {
	if !PathExists(path) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
}

func MvFile(src, dest string) error {
	_ = os.RemoveAll(dest)
	return os.Rename(src, dest)
}

func MvDir(src, dest string) error {
	_ = os.RemoveAll(dest)
	return os.Rename(src, dest)
}

func LinkFile(src, dest string) error {
	RmFile(dest)
	return os.Symlink(src, dest)
}

func RmFile(dest string) {
	_ = os.Remove(dest)
}

func RmDir(dest string) {
	_ = os.RemoveAll(dest)
}

func ChmodFile(name string, mode uint32) {
	_ = os.Chmod(name, os.FileMode(mode))
}

func IsEmptyDir(path string) bool {
	if !PathExists(path) {
		return true
	}
	if !PathIsDir(path) {
		klog.Warningf("Input path[%s] is not dir")
		return false
	}

	dir, _ := ioutil.ReadDir(path)
	return len(dir) == 0
}
