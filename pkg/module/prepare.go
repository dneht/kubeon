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

package module

import (
	"crypto/tls"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func PrepareInstall(nodes cluster.NodeList, isUpgrade bool) (err error) {
	err = prepareLocal()
	if nil != err {
		return err
	}

	err = sendPackage(nodes, isUpgrade)
	if nil != err {
		return err
	}
	err = handlePackage(nodes, isUpgrade)
	if nil != err {
		return err
	}
	return nil
}

func AfterUpgrade(node *cluster.Node, isBootstrap bool) (err error) {
	err = afterUpgrade(node)
	if nil != err {
		return err
	}
	if isBootstrap {
		client := &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
		for i := 0; i < 15; i++ {
			_, err = client.Get(node.Healthz())
			if nil == err {
				return nil
			} else {
				log.Debug("bootstrap not ready, sleep 4s")
				time.Sleep(4 * time.Second)
			}
		}
		return errors.New("check bootstrap health failed")
	}
	return nil
}

func prepareLocal() (err error) {
	current := cluster.Current()
	if nil == current {
		return errors.New("cluster not init")
	}

	current.UsePatch, err = release.PrepareLocal(cluster.CurrentResource())
	if nil != err {
		return err
	}
	release.AddLocalAutoCompletion()
	return nil
}

func sendPackage(nodes cluster.NodeList, isUpgrade bool) error {
	localRes := cluster.CurrentResource()
	var err error
	for _, node := range nodes {
		log.Infof("start copy resource to [%s]", node.Addr())
		remoteRes := node.GetResource()
		err = copyToNode(node, remoteRes.KubeletPath, localRes.KubeletPath, localRes.KubeletSum)
		if nil != err {
			log.Errorf("copy binary to [%s] failed", node.Addr())
			return err
		}
		if cluster.Current().RuntimeMode == define.ContainerdRuntime {
			err = copyToNode(node, remoteRes.ContainerdPath, localRes.ContainerdPath, localRes.ContainerdSum)
			if nil != err {
				log.Errorf("copy containerd to [%s] failed", node.Addr())
				return err
			}
		} else {
			err = copyToNode(node, remoteRes.DockerPath, localRes.DockerPath, localRes.DockerSum)
			if nil != err {
				log.Errorf("copy docker to [%s] failed", node.Addr())
				return err
			}
		}
		err = copyToNode(node, remoteRes.NetworkPath, localRes.NetworkPath, localRes.NetworkSum)
		if nil != err {
			log.Errorf("copy network to [%s] failed", node.Addr())
			return err
		}
		err = copyToNode(node, remoteRes.ImagesPath, localRes.ImagesPath, localRes.ImagesSum)
		if nil != err {
			log.Errorf("copy images to [%s] failed", node.Addr())
			return err
		}
		log.Infof("copy resource to [%s] success", node.Addr())
	}
	return nil
}

func handlePackage(nodes cluster.NodeList, upgrade bool) (err error) {
	localRes := cluster.CurrentResource()
	localConf := localRes.ClusterConf
	for _, node := range nodes {
		remoteRes := node.GetResource()
		log.Infof("start install [%s] on [%s]", define.KubeletModule, node.Addr())
		err = installOnNode(node, define.KubeletModule, remoteRes.KubeletPath)
		if nil != err {
			return err
		}
		if !upgrade {
			err = prepareScript(node)
			if nil != err {
				return err
			}
		}

		log.Infof("start install [%s] on [%s]", cluster.Current().RuntimeMode, node.Addr())
		if cluster.Current().RuntimeMode == define.ContainerdRuntime {
			err = installOnNode(node, define.ContainerdRuntime, remoteRes.ContainerdPath)
			if nil != err {
				return err
			}

			if !upgrade {
				err = enableModuleOneNow(node, define.ContainerdRuntime)
			}
			if nil != err {
				return err
			}
			err = containerdLoadImage(node)
			if nil != err {
				return err
			}
		} else {
			err = installOnNode(node, define.DockerRuntime, remoteRes.DockerPath)
			if nil != err {
				return err
			}

			if !upgrade {
				err = enableModuleOneNow(node, define.DockerRuntime)
			}
			if nil != err {
				return err
			}
			err = dockerLoadImage(node)
			if nil != err {
				return err
			}
		}
		log.Infof("start install [%s] on [%s]", define.NetworkPlugin, node.Addr())
		err = installOnNode(node, define.NetworkPlugin, remoteRes.NetworkPath)
		if nil != err {
			return err
		}
		if !upgrade {
			err = configKubeletOne(node, localConf)
		}
		if nil != err {
			return err
		}
	}
	return nil
}

func afterUpgrade(node *cluster.Node) (err error) {
	log.Infof("start reload [%s] on [%s]", cluster.Current().RuntimeMode, node.Addr())
	if release.IsUpdateCRI {
		if cluster.Current().RuntimeMode == define.ContainerdRuntime {
			err = restartModuleOne(node, define.ContainerdRuntime)
			if nil != err {
				return err
			}
		} else {
			err = restartModuleOne(node, define.DockerRuntime)
			if nil != err {
				return err
			}
		}
	}
	return restartModuleOne(node, define.KubeletModule)
}

func prepareScript(node *cluster.Node) (err error) {
	current := cluster.Current()
	installMode := "online"
	if current.IsOffline {
		installMode = "offline"
	}
	proxyMode := current.ProxyMode
	log.Infof("start prepare install on [%s], %s, proxy=%s", node.Addr(), installMode, proxyMode)
	err = node.RunCmd("bash", node.GetResource().ScriptDir+"/prepare.sh",
		"prepare", installMode, proxyMode)
	if nil != err {
		log.Errorf("prepare install on [%s] failed", node.Addr())
		return err
	}
	return nil
}

func containerdLoadImage(node *cluster.Node) (err error) {
	log.Infof("start load images on [%s]", node.Addr())
	err = node.RunCmd("ctr", "-n=k8s.io", "image", "import",
		node.GetResource().ImagesPath)
	if nil != err {
		log.Errorf("load images on [%s] failed", node.Addr())
		return err
	}
	return nil
}

func dockerLoadImage(node *cluster.Node) (err error) {
	log.Infof("start load images on [%s]", node.Addr())
	err = node.RunCmd("docker", "load", "-i",
		node.GetResource().ImagesPath)
	if nil != err {
		log.Errorf("load images on [%s] failed", node.Addr())
		return err
	}
	return nil
}
