/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package cluster

import (
	"encoding/json"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"os"
)

var runConfig *RunConfig

type RunConfig struct {
	Name  string
	Path  string
	Exist bool
	data  []byte
}

func GetConfig() *RunConfig {
	return runConfig
}

func DelConfig() {
	onutil.RmFile(runConfig.Path)
}

func InitConfig(name string) *RunConfig {
	if nil == runConfig {
		path := define.AppBaseDir + "/" + name + ".json"
		runConfig = &RunConfig{
			Name:  name,
			Path:  define.AppBaseDir + "/" + name + ".json",
			Exist: onutil.PathExists(path),
		}
	}
	return runConfig
}

func (c *RunConfig) ReadConfig() (string, error) {
	fileData, err := os.ReadFile(c.Path)
	if nil != err {
		return "", err
	}
	c.data = fileData
	return string(fileData), nil
}

func (c *RunConfig) WriteConfig() error {
	if nil == current {
		return errors.New("cluster has not been initialized, waiting...")
	}
	onutil.MkDir(define.AppBaseDir)
	clusterData, err := onutil.PrettyJson(current)
	if nil != err {
		return err
	}
	err = onutil.WriteFile(c.Path, clusterData)
	if nil != err {
		return err
	}
	c.data = clusterData
	return nil
}

func (c *RunConfig) ParseConfig() (*Cluster, error) {
	if nil == c.data {
		_, err := c.ReadConfig()
		if nil != err {
			return nil, err
		}
	}
	if nil == c.data {
		return nil, errors.New("cluster config not exist, please check process")
	}
	current = new(Cluster)
	err := json.Unmarshal(c.data, current)
	if nil != err {
		return nil, err
	}
	return current, nil
}

func (c *RunConfig) ChangeConfig() error {
	if !onutil.PathExists(current.AdminConfigPath) {
		klog.Warningf("Cluster config not created, please use: kubeon change -N %s", current.Name)
		return nil
	}

	if onutil.PathExists(define.KubernetesDefaultConfigPath) {
		onutil.RmFile(define.KubernetesDefaultConfigPath)
	} else {
		onutil.MkDir(onutil.K8sDir())
	}
	err := setLocalHost()
	if nil != err {
		return err
	}
	return onutil.LinkFile(current.AdminConfigPath, define.KubernetesDefaultConfigPath)
}
