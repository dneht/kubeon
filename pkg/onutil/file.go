/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package onutil

import (
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
		if dir, err := os.ReadDir(path); nil == err {
			for _, one := range dir {
				if !one.IsDir() {
					if info, err := one.Info(); nil == err {
						list = append(list, info)
					}
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

	dir, _ := os.ReadDir(path)
	return len(dir) == 0
}
