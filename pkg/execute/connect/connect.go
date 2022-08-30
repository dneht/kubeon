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

package connect

import (
	"github.com/dneht/kubeon/pkg/onutil"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"k8s.io/klog/v2"
	"net"
	"os"
	"sync"
	"time"
)

type AuthConfig struct {
	User       string
	Password   string
	Host       string
	Port       uint64
	PkFile     string
	PkPassword string
}

const clientTimeout = 300 * time.Second

type sshValue struct {
	client    *ssh.Client
	createdAt *time.Time
}

var sshCache = sync.Map{}

func SSHConnect(addr string) (*ssh.Session, error) {
	client, err := sshClient(addr)
	if nil != err {
		return nil, err
	}

	session, err := client.NewSession()
	if nil != err {
		return nil, err
	}
	err = session.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	})
	if nil != err {
		return nil, err
	}
	return session, nil
}

func SFTPConnect(addr string) (*sftp.Client, error) {
	client, err := sshClient(addr)
	if nil != err {
		return nil, err
	}

	connect, err := sftp.NewClient(client)
	if nil != err {
		return nil, err
	}
	return connect, nil
}

func sshClient(addr string) (*ssh.Client, error) {
	get, ok := sshCache.Load(addr)
	if ok {
		val := get.(*sshValue)
		if val.createdAt.Add(clientTimeout).Before(time.Now()) {
			return val.client, nil
		}
		sshCache.Delete(addr)
	}

	config := GetAuthConfig(addr)
	auth := sshAuthMethod(config.Password, config.PkFile, config.PkPassword)
	clientConfig := &ssh.ClientConfig{
		User:    config.User,
		Auth:    auth,
		Timeout: clientTimeout,
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if nil != err {
		return nil, err
	}
	createdAt := time.Now().Add(-10 * time.Second)
	sshCache.Store(addr, &sshValue{
		client:    client,
		createdAt: &createdAt,
	})
	return client, nil
}

func sshAuthMethod(passwd, pkFile, pkPasswd string) (auth []ssh.AuthMethod) {
	isSet := false
	if "" != pkFile && onutil.PathExists(pkFile) {
		am, err := sshPrivateKeyMethod(pkFile, pkPasswd)
		if err == nil {
			auth = append(auth, am)
			isSet = true
		}
	}
	if !isSet && "" != passwd {
		auth = append(auth, sshPasswordMethod(passwd))
		isSet = true
	}
	if !isSet {
		am, err := sshPrivateKeyMethod(onutil.Home()+"/.ssh/id_rsa", pkPasswd)
		if err == nil {
			auth = append(auth, am)
		}
	}
	return auth
}

func sshPrivateKeyMethod(pkFile, pkPassword string) (am ssh.AuthMethod, err error) {
	pkData, err := os.ReadFile(pkFile)
	if err != nil {
		klog.Errorf("[remote] Read %s file err is : %s", pkFile, err)
		os.Exit(1)
	}
	var pk ssh.Signer
	if pkPassword == "" {
		pk, err = ssh.ParsePrivateKey(pkData)
		if err != nil {
			return nil, err
		}
	} else {
		bufPwd := []byte(pkPassword)
		pk, err = ssh.ParsePrivateKeyWithPassphrase(pkData, bufPwd)
		if err != nil {
			return nil, err
		}
	}
	return ssh.PublicKeys(pk), nil
}

func sshPasswordMethod(passwd string) ssh.AuthMethod {
	return ssh.Password(passwd)
}
