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

package configset

const haproxyStaticYaml = `apiVersion: v1
kind: Pod
metadata:
  labels:
    component: local-haproxy
    tier: worker
  name: local-haproxy
  namespace: kube-system
spec:
  containers:
  - name: local-haproxy
    command:
    - local-haproxy
    {{- range .MasterHosts}}
    - --host={{.}}
    {{- end}}
    image: {{.ImageUrl}}
    imagePullPolicy: IfNotPresent
    startupProbe:
      failureThreshold: 24
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 8842
        scheme: HTTP
      initialDelaySeconds: 10
      periodSeconds: 10
      timeoutSeconds: 15
    livenessProbe:
      failureThreshold: 8
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 8842
        scheme: HTTP
      initialDelaySeconds: 10
      periodSeconds: 10
      timeoutSeconds: 15
    resources:
      requests:
        cpu: 100m
  hostNetwork: true
  priorityClassName: system-node-critical
status: {}
`
