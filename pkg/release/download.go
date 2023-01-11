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
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
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
	PauseResource      = "kubeon/install-pause"
	OfflineResource    = "kubeon/install-offline"
	DockerResource     = "kubeon/runtime-docker"
	ContainerdResource = "kubeon/runtime-containerd"
	NvidiaResource     = "kubeon/extend-nvidia"
	KataResource       = "kubeon/extend-kata"
	NetworkResource    = "kubeon/network-plugins"
	CalicoResource     = "kubeon/network-calico"
	CiliumResource     = "kubeon/network-cilium"
	ContourResource    = "kubeon/extend-contour"
	IstioResource      = "kubeon/extend-istio"
	KruiseResource     = "kubeon/extend-kruise"
)

var (
	IsUpdateRuntime = false
)

type ProcessModule struct {
	Version  string
	Resource string
	Module   string
}

func ProcessDownload(resource *ClusterResource, version, runtime, network, ingress, mirror string,
	isLocal, isBinary, isOffline, useNvidia, useKata, useKruise bool) error {
	klog.V(1).Infof("Start download binary package, version=%s, runtime=%s", version, runtime)
	if nil == resource {
		return errors.New("cluster not init")
	}
	tasks := make([]*ProcessModule, 0, 16)
	if isBinary {
		if !onutil.PathExists(resource.BinaryPath) || execute.FileSum(resource.BinaryPath) != resource.BinarySum {
			tasks = append(tasks, &ProcessModule{version, BinaryResource, define.BinaryModule})
		}
	} else {
		if !onutil.PathExists(resource.KubeletPath) || execute.FileSum(resource.KubeletPath) != resource.KubeletSum {
			tasks = append(tasks, &ProcessModule{version, DefaultResource, define.KubeletModule})

		}
	}
	if define.DockerRuntime == runtime {
		if !onutil.PathExists(resource.ContainerdPath) || execute.FileSum(resource.ContainerdPath) != resource.ContainerdSum {
			IsUpdateRuntime = true
			tasks = append(tasks, &ProcessModule{version, ContainerdResource, define.ContainerdRuntime})
		}
		if !onutil.PathExists(resource.DockerPath) || execute.FileSum(resource.DockerPath) != resource.DockerSum {
			IsUpdateRuntime = true
			tasks = append(tasks, &ProcessModule{version, DockerResource, define.DockerRuntime})
		}
	} else {
		if !onutil.PathExists(resource.ContainerdPath) || execute.FileSum(resource.ContainerdPath) != resource.ContainerdSum {
			IsUpdateRuntime = true
			tasks = append(tasks, &ProcessModule{version, ContainerdResource, define.ContainerdRuntime})
		}
	}
	if !onutil.PathExists(resource.NetworkPath) || execute.FileSum(resource.NetworkPath) != resource.NetworkSum {
		tasks = append(tasks, &ProcessModule{version, NetworkResource, define.NetworkPlugin})
	}
	if isLocal {
		if !onutil.PathExists(resource.ImagesPath) || execute.FileSum(resource.ImagesPath) != resource.ImagesSum {
			tasks = append(tasks, &ProcessModule{version, ImagesResource, define.ImagesPackage})
		}

		if extVersion, extExist := define.SupportComponentFull[version]; !extExist {
			return errors.New("extend resource not exist, please enter newer version")
		} else {
			if isOffline {
				if !onutil.PathExists(resource.OfflinePath) || execute.FileSum(resource.OfflinePath) != resource.OfflineSum {
					tasks = append(tasks, &ProcessModule{extVersion.Offline, OfflineResource, define.OfflineModule})
				}
			}
			switch network {
			case define.CalicoNetwork:
				{
					if !onutil.PathExists(resource.CalicoPath) || execute.FileSum(resource.CalicoPath) != resource.CalicoSum {
						tasks = append(tasks, &ProcessModule{extVersion.Calico, CalicoResource, define.CalicoNetwork})
					}
					break
				}
			case define.CiliumNetwork:
				{
					if !onutil.PathExists(resource.CiliumPath) || execute.FileSum(resource.CiliumPath) != resource.CiliumSum {
						tasks = append(tasks, &ProcessModule{extVersion.Cilium, CiliumResource, define.CiliumNetwork})
					}
					break
				}
			}
			switch ingress {
			case define.ContourIngress:
				{
					if !onutil.PathExists(resource.ContourPath) || execute.FileSum(resource.ContourPath) != resource.ContourSum {
						tasks = append(tasks, &ProcessModule{extVersion.Contour, ContourResource, define.ContourIngress})
					}
					break
				}
			case define.IstioIngress:
				{
					if !onutil.PathExists(resource.IstioPath) || execute.FileSum(resource.IstioPath) != resource.IstioSum {
						tasks = append(tasks, &ProcessModule{extVersion.Istio, IstioResource, define.IstioIngress})
					}
					break
				}
			}
			if useNvidia {
				if !onutil.PathExists(resource.NvidiaPath) || execute.FileSum(resource.NvidiaPath) != resource.NvidiaSum {
					tasks = append(tasks, &ProcessModule{extVersion.Nvidia, NvidiaResource, define.NvidiaRuntime})

				}
			}
			if useKata {
				if !onutil.PathExists(resource.KataPath) || execute.FileSum(resource.KataPath) != resource.KataSum {
					tasks = append(tasks, &ProcessModule{extVersion.Kata, KataResource, define.KataRuntime})
				}
			}
			if useKruise {
				if !onutil.PathExists(resource.KruisePath) || execute.FileSum(resource.KruisePath) != resource.KruiseSum {
					tasks = append(tasks, &ProcessModule{extVersion.Kruise, KruiseResource, define.KruisePlugin})
				}
			}
		}
	} else {
		if !onutil.PathExists(resource.PausePath) || execute.FileSum(resource.PausePath) != resource.PauseSum {
			tasks = append(tasks, &ProcessModule{version, PauseResource, define.PausePackage})
		}
	}

	prog := mpb.New(mpb.WithWidth(80))
	for _, task := range tasks {
		down := prog.New(0,
			mpb.BarStyle().Rbound("|"),
			mpb.PrependDecorators(
				decor.Name(task.Module, decor.WC{W: 12, C: decor.DidentRight}),
				decor.CountersKibiByte("% .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WC{W: 5}),
				decor.Name(" ] "),
			),
		)
		move := prog.New(1,
			mpb.NopStyle(),
			mpb.BarQueueAfter(down, true),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				decor.Name(task.Module, decor.WC{W: 12, C: decor.DidentRight}),
				decor.OnComplete(decor.Name("\x1b[36mmove...\x1b[0m", decor.WCSyncSpaceR), "\x1b[32mdone!\x1b[0m"),
				decor.OnComplete(decor.CountersNoUnit("%d / %d", decor.WCSyncWidth), ""),
			),
		)
		go progressImage(down, move, mirror, task.Version, task.Resource, task.Module)
	}
	prog.Wait()
	return nil
}

func progressImage(bar, move *mpb.Bar, mirror, version, image, module string) {
	temp, src, down := define.AppTmpDir+"/"+module, define.AppTmpDir+"/"+module+"/on", define.AppTmpDir+"/"+module+".tar"
	onutil.RmDir(temp)
	onutil.RmFile(down)
	onutil.MkDir(temp)

	hash, err := DownloadImage(bar, version, image, down, mirror, module)
	if nil != err {
		klog.Errorf("Download resource failed: %v", err.Error())
		os.Exit(1)
	}
	if !onutil.PathExists(down) {
		klog.Errorf("Download resource failed, please check network")
		os.Exit(1)
	}
	if err = execute.UnpackTar(down, temp); nil != err {
		klog.Errorf("Unpack resource failed, please try again")
		os.Exit(1)
	}
	if err = execute.UnpackTarGz(layerFileByBase(temp, hash), temp); nil != err {
		klog.Errorf("Unpack resource failed, please try again")
		os.Exit(1)
	}
	if !onutil.PathIsDir(src) {
		klog.Errorf("Download resource failed, please try again")
		os.Exit(1)
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
	move.SetCurrent(1)
}

func DownloadImage(bar *mpb.Bar, version, image, dest, mirror, module string) (string, error) {
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
	total, hash, err := parseImageLayer(remoteImage)
	if nil != err {
		return "", err
	}
	onutil.MkDir(define.AppTmpDir)
	bar.SetTotal(total, false)
	ticker := time.NewTicker(time.Second)
	defer func() {
		bar.SetTotal(total, true)
		ticker.Stop()
	}()
	go showProgress(bar, ticker, total, dest)
	err = tarball.WriteToFile(dest, ref, remoteImage)
	return hash, err
}

func parseImageLayer(image pkgv1.Image) (int64, string, error) {
	layers, err := image.Layers()
	if nil != err {
		return 0, "", err
	}
	size, hash, max, total := int64(0), "", int64(0), int64(0)
	for _, layer := range layers {
		size, err = layer.Size()
		if nil != err {
			return 0, "", err
		}
		if size > max {
			max = size
			get, derr := layer.Digest()
			if nil != derr {
				return 0, "", derr
			}
			hash = get.Hex
		}
		total += size
	}
	return total, hash, nil
}

func showProgress(bar *mpb.Bar, ticker *time.Ticker, total int64, dest string) {
	for range ticker.C {
		if stat, err := os.Stat(dest); nil == err {
			done := stat.Size()
			if bar.Current() != done {
				bar.SetCurrent(done)
				bar.DecoratorEwmaUpdate(time.Second)
			}
		} else {
			klog.Errorf("Check download failed: %v", err.Error())
			os.Exit(1)
		}
		if bar.Completed() {
			break
		}
	}
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
