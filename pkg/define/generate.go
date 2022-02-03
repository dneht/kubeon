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

package define

import (
	"github.com/dneht/kubeon/pkg/onutil"
	"k8s.io/klog/v2"
	"net"
)

type DefaultList struct {
	DefaultPort       uint
	DefaultUser       string
	DefaultPassword   string
	DefaultPkFile     string
	DefaultPkPassword string
}

type MasterList struct {
	MasterIPs               []net.IP
	MasterPorts             []uint
	MasterNames             []string
	MasterLabels            []string
	MasterUsers             []string
	MasterPasswords         []string
	MasterPkFiles           []string
	MasterPkPasswords       []string
	MasterDefaultPort       uint
	MasterDefaultUser       string
	MasterDefaultPassword   string
	MasterDefaultPkFile     string
	MasterDefaultPkPassword string
}

func (l MasterList) CheckMatch() bool {
	isDup := onutil.IsDuplicateInStringArr(onutil.IP2StringArr(l.MasterIPs))
	if isDup {
		klog.Error("Cluster master is duplicate")
		return false
	}
	ipSize := len(l.MasterIPs)
	if ipSize == 0 {
		klog.Error("Cluster must has one node")
		return false
	}

	return checkHostnameMatch(l.MasterNames, ipSize, "master") &&
		checkOtherMatch(l.MasterPorts, l.MasterUsers, l.MasterPkPasswords, ipSize, "master")
}

type WorkerList struct {
	WorkerIPs               []net.IP
	WorkerPorts             []uint
	WorkerNames             []string
	WorkerLabels            []string
	WorkerUsers             []string
	WorkerPasswords         []string
	WorkerPkFiles           []string
	WorkerPkPasswords       []string
	WorkerDefaultPort       uint
	WorkerDefaultUser       string
	WorkerDefaultPassword   string
	WorkerDefaultPkFile     string
	WorkerDefaultPkPassword string
}

func (l WorkerList) CheckMatch() bool {
	isDup := onutil.IsDuplicateInStringArr(onutil.IP2StringArr(l.WorkerIPs))
	if isDup {
		klog.Error("Cluster worker is duplicate")
		return false
	}

	ipSize := len(l.WorkerIPs)
	if ipSize == 0 {
		var emptyUintArr []uint
		var emptyStringArr []string
		l.WorkerPorts = emptyUintArr
		l.WorkerUsers = emptyStringArr
		l.WorkerPasswords = emptyStringArr
		l.WorkerPkFiles = emptyStringArr
		l.WorkerPkPasswords = emptyStringArr
		return true
	}

	return checkHostnameMatch(l.WorkerNames, ipSize, "worker") &&
		checkOtherMatch(l.WorkerPorts, l.WorkerUsers, l.WorkerPasswords, ipSize, "worker")
}

func checkHostnameMatch(names []string, expect int, mode string) bool {
	nameSize := len(names)
	if nameSize > 0 {
		isDup := onutil.IsDuplicateInStringArr(names)
		if isDup {
			klog.Errorf("[%s] hostnames is duplicate", mode)
			return false
		}
		if nameSize < expect {
			klog.Errorf("Number of[%s hostnames is less than expected %s", mode, expect)
			return false
		}
	}
	return true
}

func checkOtherMatch(ports []uint, users, passwords []string, expect int, mode string) bool {
	portSize := len(ports)
	if portSize > 0 && portSize < expect {
		klog.Errorf("Number of %s ports is less than expected %s", mode, expect)
		return false
	}
	userSize := len(users)
	if userSize > 0 && userSize < expect {
		klog.Errorf("Number of %s users is less than expected %s", mode, expect)
		return false
	}
	pwdSize := len(passwords)
	if pwdSize > 0 && pwdSize < expect {
		klog.Errorf("Number of %s passwords is less than expected %s", mode, expect)
		return false
	}
	return true
}
