/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
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
