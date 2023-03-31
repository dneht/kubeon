/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package action

import (
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"time"
)

func HAProxyInitWait(current *cluster.Cluster, node *cluster.Node, wait time.Duration) error {
	klog.V(1).Infof("Waiting for local-haproxy pod to become ready (timeout %s)", wait)
	if pass := waitFor(current, node, wait,
		staticPodIsReady("local-haproxy"),
	); !pass {
		return errors.New("timeout: LocalLB local-haproxy did not reach target state")
	}
	return nil
}
