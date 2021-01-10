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
	"errors"
	"github.com/dneht/kubeon/pkg/execute/connect"
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/dneht/kubeon/pkg/onutil/log"
	"github.com/pkg/sftp"
	"os"
	"path"
	"strings"
)

type RemoteCopy struct {
	node string
	src  string
	dest string
	sum  string
}

// NewRemoteCopy returns a new RemoteCopy to run copy file.
func NewRemoteCopy(node, src, dest, sum string) *RemoteCopy {
	return &RemoteCopy{
		node: node,
		src:  src,
		dest: dest,
		sum:  sum,
	}
}

// CopyTo copies the source file from local to remote.
func (c *RemoteCopy) CopyTo() error {
	log.Debugf("[remote][%s] push -- from: %s, to: %s", c.node, c.src, c.dest)
	client, err := c.getInnerCopy()
	if nil != err {
		return err
	}
	defer client.Close()

	srcFile, err := os.Open(c.src)
	if nil != err {
		log.Errorf("[remote][%s] open file [%s] err : %s", c.node, c.src, err)
		return err
	}
	defer srcFile.Close()

	destPath := path.Dir(c.dest)
	err = client.MkdirAll(destPath)
	if nil != err {
		log.Errorf("[remote][%s] create remote dir err:", c.node, destPath)
		return err
	}
	dstFile, err := client.Create(c.dest)
	if nil != err {
		log.Errorf("[remote][%s] create remote file err:", c.node, err)
		return err
	}
	defer dstFile.Close()

	sreStat, err := srcFile.Stat()
	if nil != err {
		log.Errorf("[remote][%s] get local file info err:", c.node, err)
		return err
	}
	total := sreStat.Size()
	done := int64(0)
	per := total / 4 / onutil.OneMBByte
	if per < 8 {
		per = 8
	} else if per > 160 {
		per = 160
	}
	buf := make([]byte, per*onutil.OneMBByte)
	for {
		size, ierr := srcFile.Read(buf)
		if size == 0 {
			break
		}
		if nil != ierr {
			log.Errorf("[remote][%s] read local file failed: %s", ierr)
		}
		size, ierr = dstFile.Write(buf[0:size])
		if nil != ierr {
			log.Errorf("[remote][%s] write remote file failed: %s", ierr)
		}
		done += int64(size)
		onutil.ShowProgress(total, done, "transfer")

	}
	success, localSum, remoteSum := c.isCopySuccess()
	if !success {
		log.Errorf("[remote][%s] copy file[%s](%s) to remote[%s](%s) validate cksum failed", c.node, c.src, localSum, c.dest, remoteSum)
		return errors.New("copy file validate error")
	}
	return nil
}

func (c *RemoteCopy) getInnerCopy() (*sftp.Client, error) {
	client, err := connect.SFTPConnect(c.node)
	if err != nil {
		log.Errorf("[remote][%s] ssh connect err: %s", c.node, err)
		return nil, err
	}
	return client, nil
}

func (c *RemoteCopy) isCopySuccess() (bool, string, string) {
	if "" == c.sum {
		return true, "", ""
	}

	sum, err := NewRemoteCmd(c.node, "cksum", c.dest).RunAndResult()
	if nil != err {
		log.Errorf("[remote][%s] get cksum error: %s", c.node, err)
		return false, "", ""
	}
	if "" == sum {
		return false, "", ""
	}
	sum = strings.TrimSpace(strings.Split(sum, " ")[0])
	return c.sum == sum, c.sum, sum
}
