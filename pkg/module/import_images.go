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

package module

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"k8s.io/klog/v2"
)

func ImportImages(local bool, node *cluster.Node) error {
	if local {
		if err := importImages(node, define.ImagesPackage, node.GetResource().ImagesPath); nil != err {
			return err
		}
	} else {
		if node.FileExist(node.GetResource().PausePath) {
			if err := importImages(node, define.PausePackage, node.GetResource().PausePath); nil != err {
				return err
			}
		} else {
			klog.V(1).Infof("[package] pause images not exist on [%s]", node.Addr())
		}
	}
	return nil
}

func importImages(node *cluster.Node, name, path string) (err error) {
	if cluster.Current().RuntimeMode == define.ContainerdRuntime {
		return containerdLoadImages(node, name, path)
	} else {
		return dockerLoadImages(node, name, path)
	}
}

func containerdLoadImages(node *cluster.Node, name, path string) (err error) {
	klog.V(1).Infof("[import] start [%s] images load on [%s]", name, node.Addr())
	err = node.RunCmd("ctr", "-n=k8s.io", "image", "import", path)
	if nil != err {
		klog.Errorf("[import] load [%s] images on [%s] failed", name, node.Addr())
		return err
	}
	return nil
}

func dockerLoadImages(node *cluster.Node, name, path string) (err error) {
	klog.V(1).Infof("[import] start [%s] images load on [%s]", name, node.Addr())
	err = node.RunCmd("docker", "load", "-i", path)
	if nil != err {
		klog.Errorf("[import] load [%s] images on [%s] failed", name, node.Addr())
		return err
	}
	return nil
}
