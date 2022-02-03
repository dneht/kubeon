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

package onutil

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
	"time"
)

func GetSecretSHA265() string {
	hash := sha256.New()
	hash.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func CertSHA256(data []byte) string {
	block, _ := pem.Decode(data)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		klog.Error("Failed to parse certificate: " + err.Error())
		return ""
	}
	pubData, err := x509.MarshalPKIXPublicKey(cert.PublicKey)

	hash := sha256.New()
	hash.Write(pubData)
	sum := hash.Sum(nil)
	return fmt.Sprintf("%x", sum)
}
