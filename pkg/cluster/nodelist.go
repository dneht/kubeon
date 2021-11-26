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
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

type NodeList []*Node

func (n NodeList) Len() int {
	return len(n)
}

func (n NodeList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n NodeList) Less(i, j int) bool {
	return n[i].Order < n[j].Order
}

func MergeNodeList(n1, n2 NodeList) NodeList {
	return append(n1, n2...)
}

func SortNodeList(n NodeList) NodeList {
	sort.Sort(n)
	return n
}

func EmptyNodeList() NodeList {
	return make(NodeList, 0)
}

func newNodeList(base define.DefaultList, master define.MasterList, worker define.WorkerList) (masterList NodeList, workerList NodeList, retResult bool) {
	startIdx := maxOrder() + 1
	masterSize := len(master.MasterIPs)
	masterList = make(NodeList, masterSize)
	for idx, one := range master.MasterIPs {
		node := &Node{
			IPv4:       one.String(),
			Port:       mergePort(master.MasterPorts, idx, master.MasterDefaultPort, base.DefaultPort),
			Role:       RoleControlPlane,
			User:       mergeUser(master.MasterUsers, idx, master.MasterDefaultUser, base.DefaultUser),
			Password:   mergeString(master.MasterPasswords, idx, master.MasterDefaultPassword, base.DefaultPassword),
			PkFile:     mergeString(master.MasterPkFiles, idx, master.MasterDefaultPkFile, base.DefaultPkFile),
			PkPassword: mergeString(master.MasterPkPasswords, idx, master.MasterDefaultPkPassword, base.DefaultPkPassword),
			Labels:     splitLabels(master.MasterLabels, idx),
			Order:      startIdx + uint(idx),
		}
		getHostname, result := mergeHostname(master.MasterNames, idx, node)
		if !result {
			return masterList, workerList, false
		}
		node.Hostname = getHostname
		masterList[idx] = node
	}
	workerList = make(NodeList, len(worker.WorkerIPs))
	for idx, one := range worker.WorkerIPs {
		node := &Node{
			IPv4:       one.String(),
			Port:       mergePort(worker.WorkerPorts, idx, worker.WorkerDefaultPort, base.DefaultPort),
			Role:       RoleWorker,
			User:       mergeUser(worker.WorkerUsers, idx, worker.WorkerDefaultUser, base.DefaultUser),
			Password:   mergeString(worker.WorkerPasswords, idx, worker.WorkerDefaultPassword, base.DefaultPassword),
			PkFile:     mergeString(worker.WorkerPkFiles, idx, worker.WorkerDefaultPkFile, base.DefaultPkFile),
			PkPassword: mergeString(worker.WorkerPkPasswords, idx, worker.WorkerDefaultPkPassword, base.DefaultPkPassword),
			Labels:     splitLabels(worker.WorkerLabels, idx),
			Order:      startIdx + uint(masterSize+idx),
		}
		getHostname, result := mergeHostname(worker.WorkerNames, idx, node)
		if !result {
			return masterList, workerList, false
		}
		node.Hostname = getHostname
		workerList[idx] = node
	}
	return masterList, workerList, checkHostIP(masterList, workerList) && checkHostname(masterList, workerList)
}

func checkExist(newNodes NodeList) (err error) {
	if nil == current {
		return nil
	}
	result := checkHostIP(currentNodes, newNodes) && checkHostname(currentNodes, newNodes)
	if !result {
		return errors.New("wait add node is exist")
	}
	return nil
}

func maxOrder() uint {
	startIdx := uint(0)
	if nil != current && nil != currentNodes && len(currentNodes) > 0 {
		for _, node := range currentNodes {
			if node.Order > startIdx {
				startIdx = node.Order
			}
		}
	}
	return startIdx
}

func mergePort(ports []uint, idx int, snd, fst uint) uint {
	if nil == ports || len(ports) <= idx {
		var res uint
		if fst != 22 {
			res = fst
		}
		if snd != 22 {
			res = snd
		}
		if res > 0 {
			return res
		} else {
			if snd > 0 {
				return snd
			} else {
				return fst
			}
		}
	}
	return ports[idx]
}

func mergeUser(arr []string, idx int, snd, fst string) string {
	if nil == arr || len(arr) <= idx {
		var res string
		if fst != "root" {
			res = fst
		}
		if snd != "root" {
			res = snd
		}
		if res != "" {
			return res
		} else {
			if snd != "" {
				return snd
			} else {
				return fst
			}
		}
	}
	return arr[idx]
}

func mergeString(arr []string, idx int, snd, fst string) string {
	if nil == arr || len(arr) <= idx {
		if snd != "" {
			return snd
		} else {
			return fst
		}
	}
	return arr[idx]
}

func splitLabels(arr []string, idx int) []string {
	labels := make([]string, 0)
	if nil == arr || len(arr) <= idx || "" == arr[idx] {
		return labels
	}
	list := strings.Split(arr[idx], ",")
	if len(list) == 1 {
		return append(labels, list[0])
	}
	for _, one := range list {
		labels = append(labels, one)
	}
	return labels
}

func mergeHostname(arr []string, idx int, n *Node) (string, bool) {
	n.SetConnect()
	if "root" == n.User {
		n.Home = "/root"
	} else {
		home, err := n.Command("echo ${HOME}").RunAndResult()
		if nil != err {
			return "", false
		}
		n.Home = home
	}

	rhn, err := n.RemoteHostname()
	if nil != err {
		return "", false
	}

	if nil == arr || len(arr) <= idx {
		return rhn, true
	} else {
		ghn := arr[idx]
		if rhn != ghn {
			err = n.ModifyHostname(ghn)
			if nil != err {
				return "", false
			}
		}
		return ghn, true
	}
}

func checkHostIP(ml, wl NodeList) bool {
	ms := len(ml)
	hl := make([]string, ms+len(wl))
	for idx, one := range ml {
		hl[idx] = one.IPv4
	}
	for idx, one := range wl {
		hl[ms+idx] = one.IPv4
	}
	isDup := onutil.IsDuplicateInStringArr(hl)
	if isDup {
		log.Error("cluster ip is duplicate")
		return false
	}
	return true
}

func checkHostname(ml, wl NodeList) bool {
	ms := len(ml)
	hl := make([]string, ms+len(wl))
	for idx, one := range ml {
		hl[idx] = one.Hostname
	}
	for idx, one := range wl {
		hl[ms+idx] = one.Hostname
	}
	isDup := onutil.IsDuplicateInStringArr(hl)
	if isDup {
		log.Error("cluster hostname is duplicate")
		return false
	}
	return true
}
