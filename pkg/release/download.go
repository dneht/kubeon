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

package release

import (
	"crypto/tls"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/google/go-containerregistry/pkg/name"
	pkgv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	DefaultResource    = "kubeon/install-kubeadm"
	BinaryResource     = "kubeon/install-binary"
	ImagesResource     = "kubeon/install-images"
	OfflineResource    = "kubeon/install-offline"
	DockerResource     = "kubeon/runtime-docker"
	ContainerdResource = "kubeon/runtime-containerd"
	NetworkResource    = "kubeon/network-plugins"
	NvidiaResource     = "kubeon/extend-nvidia"
	KataResource       = "kubeon/extend-kata"
	ContourResource    = "kubeon/extend-contour"
)

var (
	IsUpdateRuntime = false
)

func ProcessDownload(resource *ClusterResource, version, runtime, mirror string, isLocal, isBinary, isOffline,
	useNvidia, useKata bool, ingressMode string) (err error) {
	klog.V(1).Infof("Start download binary package, version=%s, runtime=%s", version, runtime)
	if nil == resource {
		return errors.New("cluster not initialized")
	}
	if isBinary {
		if !onutil.PathExists(resource.BinaryPath) || execute.FileSum(resource.BinaryPath) != resource.BinarySum {
			err = processImage(mirror, version, BinaryResource, define.BinaryModule)
			if nil != err {
				return err
			}
		}
	} else {
		if !onutil.PathExists(resource.KubeletPath) || execute.FileSum(resource.KubeletPath) != resource.KubeletSum {
			err = processImage(mirror, version, DefaultResource, define.KubeletModule)
			if nil != err {
				return err
			}
		}
	}
	err = downloadRuntime(resource, runtime, version, mirror)
	if nil != err {
		return err
	}
	err = downloadNetwork(resource, version, mirror)
	if nil != err {
		return err
	}
	if isLocal {
		klog.V(1).Info("Start download images package...")
		if !onutil.PathExists(resource.ImagesPath) || execute.FileSum(resource.ImagesPath) != resource.ImagesSum {
			err = processImage(mirror, version, ImagesResource, define.ImagesPackage)
			if nil != err {
				return err
			}
		}
		if isOffline {
			klog.V(1).Info("Start download offline package...")
			if !onutil.PathExists(resource.OfflinePath) || execute.FileSum(resource.OfflinePath) != resource.OfflineSum {
				err = processImage(mirror, version, OfflineResource, define.OfflineModule)
				if nil != err {
					return err
				}
			}
		}
		err = downloadExtend(resource, version, mirror, useNvidia, useKata, ingressMode)
		if nil != err {
			return err
		}
	}
	return nil
}

func downloadRuntime(localRes *ClusterResource, runtime, version, mirror string) (err error) {
	if "" == runtime || define.DockerRuntime == runtime {
		if !onutil.PathExists(localRes.ContainerdPath) || execute.FileSum(localRes.ContainerdPath) != localRes.ContainerdSum {
			IsUpdateRuntime = true
			err = processImage(mirror, version, ContainerdResource, define.ContainerdRuntime)
			if nil != err {
				return err
			}
		}
		if !onutil.PathExists(localRes.DockerPath) || execute.FileSum(localRes.DockerPath) != localRes.DockerSum {
			IsUpdateRuntime = true
			err = processImage(mirror, version, DockerResource, define.DockerRuntime)
			if nil != err {
				return err
			}
		}
	} else {
		if !onutil.PathExists(localRes.ContainerdPath) || execute.FileSum(localRes.ContainerdPath) != localRes.ContainerdSum {
			IsUpdateRuntime = true
			err = processImage(mirror, version, ContainerdResource, define.ContainerdRuntime)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

func downloadNetwork(localRes *ClusterResource, version, mirror string) (err error) {
	if !onutil.PathExists(localRes.NetworkPath) || execute.FileSum(localRes.NetworkPath) != localRes.NetworkSum {
		err = processImage(mirror, version, NetworkResource, define.NetworkPlugin)
		if nil != err {
			return err
		}
	}
	return nil
}

func downloadExtend(localRes *ClusterResource, version, mirror string, useNvidia, useKata bool, ingressMode string) (err error) {
	extVersion, ok := define.SupportComponentFull[version]
	if !ok {
		return errors.New("not found extend resource")
	}
	if useNvidia {
		klog.V(1).Info("Start download nvidia package...")
		if !onutil.PathExists(localRes.NvidiaPath) || execute.FileSum(localRes.NvidiaPath) != localRes.NvidiaSum {
			err = processImage(mirror, extVersion.Nvidia, NvidiaResource, define.NvidiaRuntime)
			if nil != err {
				return err
			}
		}
	}
	if useKata {
		klog.V(1).Info("Start download kata package...")
		if !onutil.PathExists(localRes.KataPath) || execute.FileSum(localRes.KataPath) != localRes.KataSum {
			err = processImage(mirror, extVersion.Kata, KataResource, define.KataRuntime)
			if nil != err {
				return err
			}
		}
	}
	switch ingressMode {
	case define.ContourIngress:
		{
			klog.V(1).Info("Start download contour package...")
			if !onutil.PathExists(localRes.ContourPath) || execute.FileSum(localRes.ContourPath) != localRes.ContourSum {
				err = processImage(mirror, extVersion.Contour, ContourResource, define.ContourIngress)
				if nil != err {
					return err
				}
			}
			break
		}
	}
	return nil
}

func getLocalSum(distPath, name string) string {
	localPath := distPath + "/" + name + ".sum"
	if !onutil.PathExists(localPath) {
		return ""
	}
	sum, err := ioutil.ReadFile(localPath)
	if nil != err {
		klog.V(4).Infof("Get local[%s] err is %s", localPath, err)
		return ""
	} else {
		sumStr := strings.TrimSpace(string(sum))
		klog.V(4).Infof("Get local[%s] sum is %s", localPath, sumStr)
		return sumStr
	}
}

func processImage(mirror, version, image, module string) error {
	temp, src, down := define.AppTmpDir+"/"+module, define.AppTmpDir+"/"+module+"/on", define.AppTmpDir+"/"+module+".tar"
	onutil.RmDir(temp)
	onutil.RmFile(down)
	onutil.MkDir(temp)

	if (image == ImagesResource || image == DefaultResource || image == BinaryResource) && !strings.HasPrefix(version, "v") {
		return errors.New("version format is error, like v1.x.x")
	}
	hash, err := DownloadImage(version, image, down, mirror, module)
	if nil != err {
		return err
	}
	if !onutil.PathExists(down) {
		return errors.New("download resource failed, please check network")
	}
	err = execute.UnpackTar(down, temp)
	if nil != err {
		return err
	}
	err = execute.UnpackTarGz(layerFileByBase(temp, hash), temp)
	if nil != err {
		return err
	}
	if !onutil.PathIsDir(src) {
		return errors.New("download resource failed, please retry")
	}
	onutil.MkDir(define.AppDistDir)
	for _, file := range onutil.LsDir(src) {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".tar") {
			err = onutil.MvFile(src+"/"+fileName, define.AppDistDir+"/"+fileName)
			if nil != err {
				klog.Errorf("Move file from [%s] to [%s] failed", src+"/"+fileName, define.AppDistDir+"/"+fileName)
			}
		}
	}
	onutil.RmDir(temp)
	onutil.RmFile(down)
	return err
}

func DownloadImage(version, image, dest, mirror, module string) (string, error) {
	if mirror != "" && mirror != "false" && mirror != "no" {
		image = mirror + "/" + image + ":" + version
	} else {
		image = image + ":" + version
	}
	ref, err := name.ParseReference(image, name.WeakValidation)
	if nil != err {
		return "", err
	}
	if ref.Context().RegistryStr() == name.DefaultRegistry {
		ref, err = normalizeReference(ref, image)
		if nil != err {
			return "", err
		}
	}
	registryName := ref.Context().RegistryStr()
	newReg, err := name.NewRegistry(registryName, name.WeakValidation, name.Insecure)
	if nil != err {
		return "", err
	}
	ref = setNewRegistry(ref, newReg)

	remoteImage, err := remote.Image(ref, remoteOptions()...)
	if nil != err {
		return "", err
	}
	totalSize, hash, err := parseImageLayer(remoteImage)
	if nil != err {
		return "", err
	}
	onutil.MkDir(define.AppTmpDir)
	timeTicker := time.NewTicker(time.Second * 5)
	showProgress(module, timeTicker, totalSize, dest)
	err = tarball.WriteToFile(dest, ref, remoteImage)
	timeTicker.Stop()
	return hash, err
}

func parseImageLayer(image pkgv1.Image) (int64, string, error) {
	var err error
	var layers []pkgv1.Layer
	layers, err = image.Layers()
	if nil != err {
		return 0, "", err
	}
	var size int64
	var hash string
	var max int64 = 0
	var total int64 = 0
	var get pkgv1.Hash
	for _, layer := range layers {
		size, err = layer.Size()
		if nil != err {
			return 0, "", err
		}
		if size > max {
			max = size
			get, err = layer.Digest()
			if nil != err {
				return 0, "", err
			}
			hash = get.Hex
		}
		total += size
	}
	return total, hash, nil
}

func showProgress(module string, timeTicker *time.Ticker, totalSize int64, destFile string) {
	go func(total int64, dest string) {
		for {
			<-timeTicker.C
			destStat, err := os.Stat(dest)
			if nil != err {
				klog.Errorf("[download] Get local file info err: %s", err)
				break
			} else {
				done := destStat.Size()
				onutil.ShowProgress(total, done, "download "+module)
				if done >= int64(float64(total)*0.98) {
					break
				}
			}
		}

	}(totalSize, destFile)
}

func layerFileByBase(baseDir, hash string) string {
	return baseDir + "/" + hash + ".tar.gz"
}

func normalizeReference(ref name.Reference, image string) (name.Reference, error) {
	if !strings.ContainsRune(image, '/') {
		return name.ParseReference("library/"+image, name.WeakValidation)
	}
	return ref, nil
}

func setNewRegistry(ref name.Reference, newReg name.Registry) name.Reference {
	switch r := ref.(type) {
	case name.Tag:
		r.Repository.Registry = newReg
		return r
	case name.Digest:
		r.Repository.Registry = newReg
		return r
	default:
		return ref
	}
}

func remoteOptions() []remote.Option {
	return []remote.Option{remote.WithTransport(makeTransport())}
}

func makeTransport() http.RoundTripper {
	var tr http.RoundTripper = http.DefaultTransport.(*http.Transport).Clone()
	tr.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return tr
}
