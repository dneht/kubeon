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
	"github.com/dneht/kubeon/pkg/onutil/log"
	"testing"
)

func TestMain(m *testing.M) {
	log.Init(4)
	m.Run()
}

func TestLocalCmd_Run(t *testing.T) {
	_ = NewLocalCmd("ls", "-la").Run()
}

func TestLocalCmd_RunAndResult(t *testing.T) {
	pwd, _ := NewLocalCmd("pwd").RunAndResult()
	t.Log(pwd)

	str, _ := NewLocalCmd("ls", "-la").RunAndResult()
	t.Log(str)

	str2, _ := NewLocalCmd("ls", "-la", "/").RunAndResult()
	t.Log(str2)

	str3, _ := NewLocalCmd("ls", "-la", "/tmp").RunAndResult()
	t.Log(str3)
}

func TestLocalCmd_RunWithEcho(t *testing.T) {
	_ = NewLocalCmd("ls", "-la").RunWithEcho()

	_ = NewLocalCmd("ls", "-la", "/").RunWithEcho()

	_ = NewLocalCmd("ls", "-la", "/tmp").RunWithEcho()
}

func TestLocalCmd_RunAndCapture(t *testing.T) {
	line, _ := NewLocalCmd("ls", "-la").RunAndCapture()
	for _, str := range line {
		t.Log(str)
	}

	line2, _ := NewLocalCmd("ls", "-la", "/").RunAndCapture()
	for _, str := range line2 {
		t.Log(str)
	}

	line3, _ := NewLocalCmd("ls", "-la", "/tmp").RunAndCapture()
	for _, str := range line3 {
		t.Log(str)
	}
}
