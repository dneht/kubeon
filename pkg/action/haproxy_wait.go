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

package action

import (
	"fmt"
	"github.com/dneht/kubeon/pkg/cluster"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/pkg/errors"
	"time"
)

func HAProxyInitWait(current *cluster.Cluster, node *cluster.Node, wait time.Duration) error {
	log.Infof("waiting for local-haproxy Pod to become ready (timeout %s)", wait)
	if pass := waitFor(current, node, wait,
		staticPodIsReady("local-haproxy"),
	); !pass {
		return errors.New("timeout: LocalLB local-haproxy did not reach target state")
	}
	fmt.Println()
	return nil
}
