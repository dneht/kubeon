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
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"io"
	"os"
	"os/exec"
	"strings"
)

type LocalCmd struct {
	command string
	args    []string
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
}

// NewLocalCmd returns a new LocalCmd to run a execute on local
func NewLocalCmd(command string, args ...string) *LocalCmd {
	return &LocalCmd{
		command: command,
		args:    args,
	}
}

// Run execute the inner execute on local
func (c *LocalCmd) Run() error {
	return c.runInnerCommand()
}

// RunWithEcho execute the inner execute on local and echoes the execute output to screen
func (c *LocalCmd) RunWithEcho() error {
	c.stdout = os.Stdout
	c.stderr = os.Stderr
	return c.runInnerCommand()
}

// RunAndResult executes the inner execute on local and return the output captured during execution
func (c *LocalCmd) RunAndResult() (result string, err error) {
	var buff bytes.Buffer
	c.stdout = &buff
	c.stderr = &buff
	err = c.runInnerCommand()

	lines := onutil.ReadLinesFromBuff(&buff)
	return strings.TrimSpace(lines[0]), err
}

// RunAndCapture executes the inner execute on local and return the output captured during execution
func (c *LocalCmd) RunAndCapture() (lines []string, err error) {
	var buff bytes.Buffer
	c.stdout = &buff
	c.stderr = &buff
	err = c.runInnerCommand()

	lines = onutil.ReadLinesFromBuff(&buff)
	return lines, err
}

// Stdin sets an io.Reader to be used for streaming data in input to the inner execute
func (c *LocalCmd) Stdin(in io.Reader) *LocalCmd {
	c.stdin = in
	return c
}

func (c *LocalCmd) runInnerCommand() (err error) {
	// create the commands
	cmd := exec.Command(c.command, c.args...)

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

	log.Debugf("[local] running %s", cmd.Args)
	err = cmd.Run()
	if nil != err && log.IsDebug() {
		log.Warnf("[local] running %s failed: %s", cmd.String(), err)
	}
	return err
}
