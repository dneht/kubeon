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
	"fmt"
	cluster "github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"sync"
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
				fmt.Println("Bootstrap not ready, sleep 4s")
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

func sendPackage(nodes cluster.NodeList, isUpgrade bool) (err error) {
	current := cluster.Current()
	localRes := cluster.CurrentResource()
	var wait sync.WaitGroup
	wait.Add(len(nodes))
	for _, node := range nodes {
		go func(node *cluster.Node) {
			klog.V(1).Infof("[prepare] start copy resource to [%s]", node.Addr())
			remoteRes := node.GetResource()
			err = copyToNode(node, remoteRes.KubeletPath, localRes.KubeletPath, localRes.KubeletSum)
			if nil != err {
				klog.Errorf("[prepare] copy binary to [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
			if current.RuntimeMode == define.ContainerdRuntime {
				err = copyToNode(node, remoteRes.ContainerdPath, localRes.ContainerdPath, localRes.ContainerdSum)
				if nil != err {
					klog.Errorf("[prepare] copy containerd to [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
			} else {
				err = copyToNode(node, remoteRes.DockerPath, localRes.DockerPath, localRes.DockerSum)
				if nil != err {
					klog.Errorf("[prepare] copy docker to [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
			}
			err = copyToNode(node, remoteRes.NetworkPath, localRes.NetworkPath, localRes.NetworkSum)
			if nil != err {
				klog.Errorf("[prepare] copy network to [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
			if current.IsRealLocal() {
				err = copyToNode(node, remoteRes.ImagesPath, localRes.ImagesPath, localRes.ImagesSum)
				if nil != err {
					klog.Errorf("[prepare] copy images to [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
				if current.IsOffline {
					err = copyToNode(node, remoteRes.OfflinePath, localRes.OfflinePath, localRes.OfflineSum)
					if nil != err {
						klog.Errorf("[prepare] copy offline to [%s] failed: %v", node.Addr(), err)
						os.Exit(1)
					}
				}
				if current.UseNvidia && node.HasNvidia {
					err = copyToNode(node, remoteRes.NvidiaPath, localRes.NvidiaPath, localRes.NvidiaSum)
					if nil != err {
						klog.Errorf("[prepare] copy nvidia to [%s] failed: %v", node.Addr(), err)
						os.Exit(1)
					}
				}
				if current.UseKata {
					err = copyToNode(node, remoteRes.KataPath, localRes.KataPath, localRes.KataSum)
					if nil != err {
						klog.Errorf("[prepare] copy kata to [%s] failed: %v", node.Addr(), err)
						os.Exit(1)
					}
				}
				switch current.IngressMode {
				case define.ContourIngress:
					{
						err = copyToNode(node, remoteRes.ContourPath, localRes.ContourPath, localRes.ContourSum)
						if nil != err {
							klog.Errorf("[prepare] copy contour to [%s] failed: %v", node.Addr(), err)
							os.Exit(1)
						}
						break
					}
				}
			}
			klog.V(1).Infof("[prepare] copy resource to [%s] success", node.Addr())
			wait.Done()
		}(node)
	}
	wait.Wait()
	return err
}

func handlePackage(nodes cluster.NodeList, upgrade bool) (err error) {
	current := cluster.Current()
	localRes := cluster.CurrentResource()
	localConf := localRes.ClusterConf
	var wait sync.WaitGroup
	wait.Add(len(nodes))
	for _, node := range nodes {
		go func(node *cluster.Node) {
			remoteRes := node.GetResource()
			klog.V(1).Infof("[package] Start install [%s] on [%s]", define.KubeletModule, node.Addr())
			_, err = installOnNode(node, define.KubeletModule, remoteRes.KubeletPath)
			if nil != err {
				klog.Errorf("[package] install kubelet on [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
			if !upgrade {
				err = prepareScript(node)
				if nil != err {
					klog.Errorf("[package] prepare script on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
			}

			nowRuntimeMode := cluster.Current().RuntimeMode
			klog.V(1).Infof("[package] start install [%s] on [%s]", nowRuntimeMode, node.Addr())
			if nowRuntimeMode == define.ContainerdRuntime {
				_, err = installOnNode(node, define.ContainerdRuntime, remoteRes.ContainerdPath)
				if nil != err {
					klog.Errorf("[package] install containerd on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
			} else {
				_, err = installOnNode(node, define.DockerRuntime, remoteRes.DockerPath)
				if nil != err {
					klog.Errorf("[package] install docker on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
			}
			if !upgrade {
				err = enableModuleOneNow(node, nowRuntimeMode)
			}
			if nil != err {
				klog.Errorf("[package] enable runtime on [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
			if current.IsRealLocal() {
				err = importImage(node, remoteRes.ImagesPath)
				if nil != err {
					klog.Errorf("[package] import base image on [%s] failed: %v", node.Addr(), err)
					os.Exit(1)
				}
			} else {
				if node.FileExist(remoteRes.PausePath) {
					err = importImage(node, remoteRes.PausePath)
					if nil != err {
						klog.Errorf("[package] import pause image on [%s] failed: %v", node.Addr(), err)
						os.Exit(1)
					}
				}
			}
			klog.V(1).Infof("[package] start install [%s] on [%s]", define.NetworkPlugin, node.Addr())
			_, err = installOnNode(node, define.NetworkPlugin, remoteRes.NetworkPath)
			if nil != err {
				klog.Errorf("[package] install cni on [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
			if current.IsRealLocal() {
				if current.UseNvidia && node.HasNvidia {
					err = importOnNode(node, define.NvidiaRuntime, remoteRes.NvidiaPath)
					if nil != err {
						klog.Errorf("[package] import nvidia image on [%s] failed: %v", node.Addr(), err)
						os.Exit(1)
					}
					setupNvidia(node, nowRuntimeMode)
				}
				if current.UseKata {
					err = importOnNode(node, define.KataRuntime, remoteRes.KataPath)
					if nil != err {
						klog.Errorf("[package] import kata image on [%s] failed: %v", node.Addr(), err)
						os.Exit(1)
					}
				}
				switch current.IngressMode {
				case define.ContourIngress:
					{
						err = importOnNode(node, define.ContourIngress, remoteRes.ContourPath)
						if nil != err {
							klog.Errorf("[package] import contour image on [%s] failed: %v", node.Addr(), err)
							os.Exit(1)
						}
						break
					}
				}
			} else {
				if current.UseNvidia && node.HasNvidia {
					setupNvidia(node, nowRuntimeMode)
				}
			}
			if !upgrade {
				err = configKubeletOne(node, localConf)
			}
			if nil != err {
				klog.Errorf("[package] enable kubelet on [%s] failed: %v", node.Addr(), err)
				os.Exit(1)
			}
			wait.Done()
		}(node)
	}
	wait.Wait()
	return nil
}

func prepareScript(node *cluster.Node) (err error) {
	current := cluster.Current()
	installMode := "online"
	if current.IsOffline {
		installMode = "offline"
	}
	proxyMode := current.ProxyMode
	klog.V(1).Infof("Start prepare install on [%s], %s, proxy=%s", node.Addr(), installMode, proxyMode)
	err = node.RunCmd("bash", node.GetResource().ScriptDir+"/prepare.sh",
		"prepare", installMode, proxyMode)
	if nil != err {
		klog.Errorf("Prepare install on [%s] failed", node.Addr())
		return err
	}
	if current.UseNvidia && node.HasNvidia {
		klog.V(1).Infof("Start discover nvidia on [%s]", node.Addr())
		err = node.RunCmd("bash", node.GetResource().ScriptDir+"/discover.sh",
			"nvidia", "no", installMode)
		if nil != err {
			klog.Errorf("Discover nvidia on [%s] failed", node.Addr())
			return err
		}
	}
	return nil
}

func afterUpgrade(node *cluster.Node) (err error) {
	klog.V(1).Infof("Start reload [%s] on [%s]", cluster.Current().RuntimeMode, node.Addr())
	if release.IsUpdateRuntime {
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

func setupNvidia(node *cluster.Node, nowRuntimeMode string) {
	var err error
	err = node.RunCmd("sed", "-i", "-E", "\"s#BinaryName\\s+=\\s+\\\"[0-9a-zA-Z\\-]+\\\"#BinaryName = \\\"nvidia-container-runtime\\\"#g\"", "/etc/containerd/config.toml")
	if nil != err {
		klog.Errorf("[package] import nvidia image on [%s] failed: %v", node.Addr(), err)
		os.Exit(1)
	}
	err = restartModuleOne(node, nowRuntimeMode)
	if nil != err {
		klog.Errorf("[package] restart runtime on [%s] failed: %v", node.Addr(), err)
		os.Exit(1)
	}
}

func importImage(node *cluster.Node, path string) (err error) {
	if cluster.Current().RuntimeMode == define.ContainerdRuntime {
		return containerdLoadImage(node, path)
	} else {
		return dockerLoadImage(node, path)
	}
}

func containerdLoadImage(node *cluster.Node, path string) (err error) {
	klog.V(1).Infof("Start load images on [%s]", node.Addr())
	err = node.RunCmd("ctr", "-n=k8s.io", "image", "import", path)
	if nil != err {
		klog.Errorf("[import] load images on [%s] failed", node.Addr())
		return err
	}
	return nil
}

func dockerLoadImage(node *cluster.Node, path string) (err error) {
	klog.V(1).Infof("Start load images on [%s]", node.Addr())
	err = node.RunCmd("docker", "load", "-i", path)
	if nil != err {
		klog.Errorf("[import] load images on [%s] failed", node.Addr())
		return err
	}
	return nil
}
