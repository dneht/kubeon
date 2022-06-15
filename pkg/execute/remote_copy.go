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
	"github.com/pkg/sftp"
	"github.com/vbauerster/mpb/v7"
	"io"
	"k8s.io/klog/v2"
	"os"
	"path"
	"strings"
)

type RemoteCopy struct {
	node string
	src  string
	dest string
	sum  string
	bar  *mpb.Bar
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

func (c *RemoteCopy) Node() string {
	return c.node
}

func (c *RemoteCopy) UseBar(bar *mpb.Bar) {
	c.bar = bar
}

// CopyTo copies the source file from local to remote.
func (c *RemoteCopy) CopyTo() error {
	client, err := c.getInnerCopy()
	if nil != err {
		return err
	}
	defer client.Close()

	srcFile, err := os.Open(c.src)
	if nil != err {
		klog.Errorf("[remote] [%s] local file [%s] err", c.node, c.src, err)
		return err
	}
	defer srcFile.Close()

	destPath := path.Dir(c.dest)
	err = client.MkdirAll(destPath)
	if nil != err {
		klog.Errorf("[remote] [%s] create remote dir err: %s", c.node, destPath)
		return err
	}
	dstFile, err := client.Create(c.dest)
	if nil != err {
		klog.Errorf("[remote] [%s] create remote file err", c.node, err)
		return err
	}
	defer dstFile.Close()

	if nil == c.bar {
		klog.V(4).Infof("[remote] [%s] push -- from: %s, to: %s", c.node, c.src, c.dest)
		_, err = io.Copy(dstFile, srcFile)
		if nil != err {
			klog.Errorf("[remote] [%s] copy file to remote failed", err)
			return err
		}
	} else {
		srcStat, serr := srcFile.Stat()
		if nil != serr {
			klog.Errorf("[remote] [%s] get local file info err", c.node, err)
			return serr
		}
		c.bar.SetTotal(srcStat.Size(), false)
		proxyReader := c.bar.ProxyReader(srcFile)
		defer func() {
			_ = proxyReader.Close()
			c.bar.SetTotal(srcStat.Size(), true)
		}()

		_, err = io.Copy(dstFile, proxyReader)
		if nil != err {
			klog.Errorf("[remote] [%s] copy file to remote failed", err)
			return err
		}
	}

	success, localSum, remoteSum := c.isCopySuccess()
	if !success {
		klog.Errorf("[remote] [%s] copy file[%s](%s) to remote[%s](%s) validate cksum failed", c.node, c.src, localSum, c.dest, remoteSum)
		return errors.New("checksum mismatch on copied file")
	}
	return nil
}

func (c *RemoteCopy) getInnerCopy() (*sftp.Client, error) {
	client, err := connect.SFTPConnect(c.node)
	if err != nil {
		klog.Errorf("[remote] [%s] connect remote err: %s", c.node, err)
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
		klog.Errorf("[remote] [%s] get cksum error", c.node, err)
		return false, "", ""
	}
	if "" == sum {
		return false, "", ""
	}
	sum = strings.TrimSpace(strings.Split(sum, " ")[0])
	return c.sum == sum, c.sum, sum
}
