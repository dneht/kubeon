/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package cluster

import (
	"github.com/dneht/kubeon/pkg/onutil/log"
	"testing"
)

func TestMain(m *testing.M) {
	log.Init(6)
	m.Run()
}

func TestNodeListSort_Run(t *testing.T) {
	list := NodeList{
		&Node{
			IPv4:  "1.1.1.1",
			Order: 2,
		},
		&Node{
			IPv4:  "9.1.1.1",
			Order: 9,
		},
		&Node{
			IPv4:  "0.1.1.1",
			Order: 0,
		},
		&Node{
			IPv4:  "3.1.1.1",
			Order: 3,
		},
	}
	for _, n := range SortNodeList(list) {
		t.Log(n.IPv4)
	}
}
