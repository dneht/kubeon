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

package execute

import (
	"github.com/dneht/kubeon/pkg/execute/connect"
	"testing"
)

const node = "192.168.2.11:22"

func init()  {
	_, _, _ = connect.SetAuthConfig(&connect.AuthConfig{
		User:     "root",
		Host:     "192.168.2.11",
		Port:     22,
	})
}

func TestRemoteCmd_Run(t *testing.T) {
	_ = NewRemoteCmd(node, "ls", "-la").Run()
}

func TestRemoteCmd_RunAndResult(t *testing.T) {
	str, _ := NewRemoteCmd(node, "ls", "-la").RunAndResult()
	t.Log(str)

	str2, _ := NewRemoteCmd(node, "ls", "-la", "/").RunAndResult()
	t.Log(str2)

	str3, _ := NewRemoteCmd(node, "ls", "-la", "/tmp").RunAndResult()
	t.Log(str3)
}

func TestRemoteCmd_RunAndCapture(t *testing.T) {
	echo, _ := NewRemoteCmd(node, "echo", "${PATH}", "&&", "echo", "${HOME}").RunAndCapture()
	for _, str := range echo {
		t.Log(str)
	}

	line, _ := NewRemoteCmd(node, "ls", "-la").RunAndCapture()
	for _, str := range line {
		t.Log(str)
	}

	line2, _ := NewRemoteCmd(node, "ls", "-la", "/").RunAndCapture()
	for _, str := range line2 {
		t.Log(str)
	}

	line3, _ := NewRemoteCmd(node, "ls", "-la", "/tmp").RunAndCapture()
	for _, str := range line3 {
		t.Log(str)
	}
}
