/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/define"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"strings"
)

var (
	etcdCertArgsNew = []string{"--cacert=/etc/kubernetes/pki/etcd/ca.crt", "--cert=/etc/kubernetes/pki/etcd/peer.crt", "--key=/etc/kubernetes/pki/etcd/peer.key"}
	etcdCertArgsOld = []string{"--ca-file=/etc/kubernetes/pki/etcd/ca.crt", "--cert-file=/etc/kubernetes/pki/etcd/peer.crt", "--key-file=/etc/kubernetes/pki/etcd/peer.key"}
)

func Etcdctl(cmds []string) error {
	boot := cluster.BootstrapNode()
	args := buildEtcdctlArgs(boot)
	// Append version specific etcdctl certificate flags
	err := appendEtcdctlCertArgs(EtcdVersion(), &args)
	if nil != err {
		return err
	}
	args = append(args, cmds...)
	return boot.Command("kubectl", args...).RunWithEcho()
}

func EtcdVersion() string {
	boot := cluster.BootstrapNode()
	args := append(buildEtcdctlArgs(boot), "version")
	lines, err := boot.Command("kubectl", args...).RunAndCapture()
	if err != nil {
		klog.Warningf("Get etcdctl version error: %s", err)
		return define.ETCD_3_4_0.Full
	}
	version, err := parseEtcdctlVersion(lines)
	if err != nil {
		klog.Warningf("Parse etcdctl version error: %s", err)
		return define.ETCD_3_4_0.Full
	}
	return version
}

func EtcdMemberList() ([]string, error) {
	boot := cluster.BootstrapNode()
	args := buildEtcdctlArgs(boot)
	// Append version specific etcdctl certificate flags
	err := appendEtcdctlCertArgs(EtcdVersion(), &args)
	if nil != err {
		return nil, err
	}
	args = append(args, "member", "list")
	return boot.Command("kubectl", args...).RunAndCapture()
}

func EtcdMemberRemove(name string) error {
	boot := cluster.BootstrapNode()
	args := buildEtcdctlArgs(boot)
	// Append version specific etcdctl certificate flags
	err := appendEtcdctlCertArgs(EtcdVersion(), &args)
	if nil != err {
		return err
	}
	args = append(args, "member", "list")
	lines, err := boot.Command("kubectl", args...).RunAndCapture()
	if nil != err {
		return err
	}
	for _, line := range lines {
		if strings.Contains(line, name) {
			id := strings.Split(line, ",")[0]
			removeArgs := append(args, "member", "remove", id)
			return boot.Command("kubectl", removeArgs...).RunWithEcho()
		}
	}
	return nil
}

func EtcdSnapshotSave(file string) error {
	return etcdSnapshotOpera("snapshot", "save", file)
}

func EtcdSnapshotCheck(file string) error {
	return etcdSnapshotOpera("snapshot", "--write-out=table", "status", file)
}

func etcdSnapshotOpera(input ...string) error {
	boot := cluster.BootstrapNode()
	args := buildEtcdctlArgs(boot)
	// Append version specific etcdctl certificate flags
	err := appendEtcdctlCertArgs(EtcdVersion(), &args)
	if nil != err {
		return err
	}
	args = append(args, input...)
	return boot.Command("kubectl", args...).RunWithEcho()
}

func buildEtcdctlArgs(boot *cluster.Node) []string {
	// NB. before v1.13 local etcd is listening on localhost only; after v1.13
	// local etcd is listening on localhost and on the advertise address; we are
	// using localhost to accommodate both the use cases
	return []string{
		"--kubeconfig=/etc/kubernetes/admin.conf", "exec", "--namespace=kube-system", fmt.Sprintf("etcd-%s", boot.Hostname),
		"--", "etcdctl",
	}
}

// parseEtcdctlVersion takes the output lines of 'etcdctl version' and returns the version
func parseEtcdctlVersion(lines []string) (string, error) {
	if len(lines) < 1 {
		return "", errors.New("expected at least one line from the output of 'etcdctl version'")
	}
	elements := strings.Split(lines[0], ":")
	if len(elements) != 2 {
		return "", errors.New("expected ':' on the first line of 'etcdctl version'")
	}
	return strings.TrimSpace(elements[1]), nil
}

// appendEtcdctlCertArgs takes an etcd version and appends etcdctl certificate arguments
// to a existing list of arguments based on the version
func appendEtcdctlCertArgs(etcdVersion string, etcdArgs *[]string) error {
	version, err := define.NewStdVersion(etcdVersion)
	if err != nil {
		return errors.Wrap(err, "cannot parse etcd version")
	}

	// Before 3.4.0, etcdctl was using --ca-file, --cert-file, --key-file flags; in newer etcdctl releases those flags are renamed
	if version.LessThen(define.ETCD_3_4_0) {
		*etcdArgs = append(*etcdArgs, etcdCertArgsOld...)
	} else {
		*etcdArgs = append(*etcdArgs, etcdCertArgsNew...)
	}
	return nil
}
