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

package cluster

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/execute"
	"github.com/dneht/kubeon/pkg/onutil"
	"k8s.io/klog/v2"
	"sort"
	"strings"
)

func InitHost() (err error) {
	all := CurrentNodes()
	kv := make(map[string]string, len(all)+1)
	boot := BootstrapNode()
	kv[current.LbDomain] = boot.IPv4
	for _, node := range all {
		kv[node.Hostname] = node.IPv4
	}
	hosts := getAddHosts(kv)
	for _, node := range all {
		err = node.RunCmd(hosts)
		if nil != err {
			klog.Warningf("Init hosts error on[%s]: %s", node.Addr(), err)
		}
	}
	return setLocalHost()
}

func UpdateHost() (err error) {
	all := CurrentNodes()
	kv := make(map[string]string, 1)
	boot := BootstrapNode()
	for _, node := range all {
		if node.IsControlPlane() {
			kv[current.LbDomain] = node.IPv4
			hosts := getAddHosts(kv)
			err = node.RunCmd(hosts)
			if nil != err {
				klog.Warningf("Init hosts error on[%s]: %s", node.Addr(), err)
			}
		} else if node.IsWorker() {
			if current.IsMultiMaster() {
				kv[current.LbDomain] = current.LbIP
			} else {
				kv[current.LbDomain] = boot.IPv4
			}
			hosts := getAddHosts(kv)
			err = node.RunCmd(hosts)
			if nil != err {
				klog.Warningf("Update hosts error on[%s]: %s", node.Addr(), err)
			}
		}
	}
	return setLocalHost()
}

func DeleteHost(delNodes NodeList) (err error) {
	all := CurrentNodes()
	kv := make(map[string]string, len(delNodes))
	for _, node := range delNodes {
		kv[node.Hostname] = node.IPv4
	}
	hosts := getDelHosts(kv)
	for _, node := range all {
		err = node.RunCmd(hosts)
		if nil != err {
			klog.Warningf("Delete hosts error on[%s]: %s", node.Addr(), err)
		}
	}
	return setLocalHost()
}

func getAddHosts(hosts map[string]string) (run string) {
	var sb strings.Builder
	keys := make([]string, 0, len(hosts))
	for key := range hosts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		sb.WriteString("sed -i -E '/^[0-9a-f.:]+\\s+")
		sb.WriteString(key)
		sb.WriteString(".*$'/d /etc/hosts")
		sb.WriteString(" && ")
		sb.WriteString("echo '")
		sb.WriteString(hosts[key])
		sb.WriteString("  ")
		sb.WriteString(key)
		sb.WriteString("' >>/etc/hosts")
		sb.WriteString(" && ")
	}
	run = sb.String()
	return run[0 : len(run)-4]
}

func getDelHosts(hosts map[string]string) (run string) {
	var sb strings.Builder
	for key, _ := range hosts {
		sb.WriteString("sed -i -E '/^[0-9a-f.:]+\\s+")
		sb.WriteString(key)
		sb.WriteString(".*$'/d /etc/hosts")
		sb.WriteString(" && ")
	}
	run = sb.String()
	return run[0 : len(run)-4]
}

func setLocalHost() (err error) {
	if !onutil.IsLocalIPv4InCluster(current.AllIPs()) {
		err = execute.NewLocalCmd("sh", "-c",
			fmt.Sprintf("sed -i -E '/^[0-9a-f.:]+\\s+%s.*$'/d /etc/hosts && echo '%s  %s' >> /etc/hosts",
				current.LbDomain, BootstrapNode().IPv4, current.LbDomain)).Run()
		if nil != err {
			klog.Warningf("Set local lb domain ip failed: %v", err)
		}
	}
	return nil
}

func resetLocalHost(node *Node) {
	if onutil.IsLocalIPv4(node.IPv4) {
		err := execute.NewLocalCmd("sh", "-c",
			fmt.Sprintf("sed -i -E '/^[0-9a-f.:]+\\s+%s.*$'/d /etc/hosts && echo '%s  %s' >> /etc/hosts",
				current.LbDomain, BootstrapNode().IPv4, current.LbDomain)).Run()
		if nil != err {
			klog.Warningf("Reset local lb domain ip failed: %v", err)
		}
	}
}
