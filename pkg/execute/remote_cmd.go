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
	"bytes"
	"github.com/dneht/kubeon/pkg/execute/connect"
	"github.com/dneht/kubeon/pkg/onutil"
	"io"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

type RemoteCmd struct {
	node    string
	command string
	args    []string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

// NewRemoteCmd returns a new RemoteCmd to run a execute on remote
func NewRemoteCmd(node, command string, args ...string) *RemoteCmd {
	return &RemoteCmd{
		node:    node,
		command: command,
		args:    args,
	}
}

// Run execute the inner execute on remote
func (c *RemoteCmd) Run() error {
	return c.runInnerCommand()
}

// RunAndResult executes the inner execute on remote and return the output captured during execution
func (c *RemoteCmd) RunAndResult() (result string, err error) {
	var normal bytes.Buffer
	c.stdout = &normal
	var wrong bytes.Buffer
	c.stderr = &wrong
	err = c.runInnerCommand()

	var lines []string
	if nil != err {
		lines = onutil.ReadLinesFromBuff(&wrong)
		result = strings.Join(lines, "\n")
		return result, err
	} else {
		if normal.Len() == 0 {
			return "", err
		}
		lines = onutil.ReadLinesFromBuff(&normal)
		return strings.TrimSpace(lines[0]), err
	}
}

// RunWithEcho execute the inner execute on remote and echoes the execute output to screen
func (c *RemoteCmd) RunWithEcho() error {
	c.stdout = os.Stderr
	c.stderr = os.Stdout
	return c.runInnerCommand()
}

// RunAndCapture executes the inner execute on remote and return the output captured during execution
func (c *RemoteCmd) RunAndCapture() (lines []string, err error) {
	var normal bytes.Buffer
	c.stdout = &normal
	var wrong bytes.Buffer
	c.stderr = &wrong
	err = c.runInnerCommand()

	if nil != err {
		lines = onutil.ReadLinesFromBuff(&wrong)
		return lines, err
	} else {
		if normal.Len() == 0 {
			return []string{}, err
		}
		lines = onutil.ReadLinesFromBuff(&normal)
		return lines, nil
	}
}

// Stdin sets an io.Reader to be used for streaming data in input to the inner execute
func (c *RemoteCmd) Stdin(in io.Reader) *RemoteCmd {
	c.stdin = in
	return c
}

func (c *RemoteCmd) runInnerCommand() error {
	cmd, err := connect.SSHConnect(c.node)
	if nil != err {
		return err
	}
	defer cmd.Close()

	// redirects flows if requested
	if nil != c.stdin {
		cmd.Stdin = c.stdin
	}
	if nil != c.stdout {
		cmd.Stdout = c.stdout
	}
	if nil != c.stderr {
		cmd.Stderr = c.stderr
	}

	run := c.command + " " + strings.Join(c.args, " ")
	klog.V(6).Infof("[remote] running [%s] on [%s]", run, c.node)
	err = cmd.Run(run)
	if nil != err {
		klog.Errorf("[remote] running [%s] on [%s] failed: %s", run, c.node, err)
	}
	return err
}
