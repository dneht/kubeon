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

package cluster

import (
	"github.com/dneht/kubeon/pkg/define"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/execute/connect"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/dneht/kubeon/pkg/release"
	"github.com/vbauerster/mpb/v7"
	"k8s.io/klog/v2"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	addr       string
	IPv4       string   `json:"ip"`
	IPv6       string   `json:"ipv6"`
	Port       uint     `json:"port"`
	Hostname   string   `json:"hostname"`
	Labels     []string `json:"labels"`
	Home       string   `json:"home"`
	Role       NodeRole `json:"role"`
	User       string   `json:"user"`
	Password   string   `json:"password"`
	Link       string   `json:"link"`
	PkFile     string   `json:"pkFile"`
	PkPassword string   `json:"pkPassword"`
	Order      uint     `json:"order"`
	Status     string   `json:"status"`
	HasNvidia  bool     `json:"hasNvidia"`
	resource   *release.ClusterRemoteResource
}

type NodeRole string

const (
	RoleControlPlane NodeRole = "cp"
	RoleWorker       NodeRole = "worker"
)

func (n *Node) Addr() string {
	if "" != n.addr {
		return n.addr
	}
	if "" == n.IPv4 {
		n.addr = n.IPv6 + ":" + strconv.FormatUint(uint64(n.Port), 10)
	} else {
		n.addr = n.IPv4 + ":" + strconv.FormatUint(uint64(n.Port), 10)
	}
	return n.addr
}

func (n *Node) IsBootstrap() bool {
	currentNodes := CurrentNodes()
	return n.Order == currentNodes[0].Order
}

func (n *Node) IsControlPlane() bool {
	return n.Role == RoleControlPlane
}

func (n *Node) IsWorker() bool {
	return n.Role == RoleWorker
}

func (n *Node) Healthz() string {
	return "https://" + n.IPv4 + ":" + strconv.FormatInt(int64(define.DefaultClusterAPIPort), 10) + "/healthz"
}

func (n *Node) RunCmd(cmd string, args ...string) error {
	run := execute.NewRemoteCmd(n.Addr(), cmd, args...)
	if log.IsDebug() {
		return run.RunWithEcho()
	} else {
		return run.Run()
	}
}

func (n *Node) Command(cmd string, args ...string) *execute.RemoteCmd {
	return execute.NewRemoteCmd(n.Addr(), cmd, args...)
}

func (n *Node) CopyTo(src, dest string) error {
	return execute.NewRemoteCopy(n.Addr(), src, dest, "").CopyTo()
}

func (n *Node) CopyToWithSum(src, dest, sum string) error {
	return execute.NewRemoteCopy(n.Addr(), src, dest, sum).CopyTo()
}

func (n *Node) CopyToWithBar(src, dest, sum string, bar *mpb.Bar) {
	rc := execute.NewRemoteCopy(n.Addr(), src, dest, sum)
	if nil != bar {
		rc.UseBar(bar)
	}
	err := rc.CopyTo()
	if nil != err {
		klog.Errorf("copy to remote failed: %v", err)
		os.Exit(1)
	}
}

func (n *Node) Chmod(mode, path string) error {
	return n.RunCmd("chmod", mode, path)
}

func (n *Node) Rm(path string) error {
	return n.RunCmd("rm", "-rf", path)
}

func (n *Node) MkDir(path string) error {
	return n.RunCmd("mkdir", "-p", path)
}

func (n *Node) DirExist(path string) bool {
	return testIsExist(n, path, "-d")
}

func (n *Node) FileExist(path string) bool {
	return testIsExist(n, path, "-f")
}

func testIsExist(n *Node, path, flag string) bool {
	result, err := n.Command("if", "test",
		flag, path, ";then", "echo", "yes;", "fi").RunAndResult()
	if nil == err && "yes" == result {
		return true
	} else {
		return false
	}
}

func (n *Node) FileSum(path string) string {
	result, err := n.Command("if", "test",
		"-f", path, ";then", "cksum", path, ";", "fi").RunAndResult()
	if nil == err && len(result) >= 4 {
		return strings.TrimSpace(strings.Split(result, " ")[0])
	} else {
		return ""
	}
}

func (n *Node) CheckNvidia() bool {
	return checkDevice(n, "nvidia")
}

func checkDevice(n *Node, device string) bool {
	result, err := n.Command("for", "name", "in", "/dev/"+device+"*;", "do", "if", "test",
		"-c", "${name}", ";then", "echo", "yes;", "break;", "fi", "done").RunAndResult()
	if nil == err && "yes" == result {
		return true
	} else {
		return false
	}
}

func (n *Node) RemoteHostname() (string, error) {
	result, err := n.Command("hostname").RunAndResult()
	if nil != err {
		klog.Warningf("Get remote[%s] hostname error: %s", n.Addr(), err)
		return "", err
	}
	return result, nil
}

func (n *Node) ModifyHostname(hostname string) error {
	err := n.Command("hostnamectl",
		"--pretty", "--static", "--transient",
		"set-hostname", hostname).Run()
	if nil != err {
		klog.Warningf("Set remote[%s] hostname error: %s", n.Addr(), err)
		return err
	}
	return nil
}

func (n *Node) KubeVersion() string {
	result, err := n.Command("cat", define.AppBaseDir+"/version/k8s").RunAndResult()
	if nil != err {
		klog.Warningf("Get k8s version error: %s", err)
		return ""
	}
	return result
}

func (n *Node) CRIVersion() string {
	result, err := n.Command("cat", define.AppBaseDir+"/version/cri").RunAndResult()
	if nil != err {
		klog.Warningf("Get cri version error: %s", err)
		return ""
	}
	return result
}

func (n *Node) CNIVersion() string {
	result, err := n.Command("cat", define.AppBaseDir+"/version/cni").RunAndResult()
	if nil != err {
		klog.Warningf("Get cni version error: %s", err)
		return ""
	}
	return result
}

func (n *Node) SetConnect() {
	_, _, err := connect.SetAuthConfig(&connect.AuthConfig{
		User:       n.User,
		Password:   n.Password,
		Host:       n.IPv4,
		Port:       uint64(n.Port),
		PkFile:     n.PkFile,
		PkPassword: n.PkPassword,
	})
	if nil != err {
		klog.Errorf("Get[%s]:[%s] ssh connect failed", n.IPv4, n.Port)
		os.Exit(1)
	}
}

func (n *Node) GetResource() *release.ClusterRemoteResource {
	if nil == n.resource {
		n.resource = release.RemoteResource(n.Home, current.RuntimeMode)
	}
	return n.resource
}

func (n *Node) LocalConfigPath() string {
	return current.LocalResource.ClusterConf.KubeadmInitDir + n.Hostname + ".yaml"
}
