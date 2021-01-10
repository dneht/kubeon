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

package configset

const kubeletYaml = `apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
address: 0.0.0.0
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 2m0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 5m0s
    cacheUnauthorizedTTL: 30s
cgroupDriver: systemd
cgroupsPerQOS: true
clusterDomain: {{.ClusterDnsDomain}}
configMapAndSecretChangeDetectionStrategy: Watch
containerLogMaxFiles: 5
containerLogMaxSize: 10Mi
enforceNodeAllocatable:
  - pods
eventBurst: 10
eventRecordQPS: 5
evictionHard:
  imagefs.available: 15%
  memory.available: 300Mi
  nodefs.available: 10%
  nodefs.inodesFree: 5%
evictionPressureTransitionPeriod: 5m0s
failSwapOn: true
hairpinMode: promiscuous-bridge
maxOpenFiles: 1000000
maxPods: {{.ClusterMaxPods}}
nodeLeaseDurationSeconds: 40
nodeStatusReportFrequency: 1m0s
nodeStatusUpdateFrequency: 10s
oomScoreAdj: -999
podPidsLimit: -1
rotateCertificates: true
runtimeRequestTimeout: 2m0s
serializeImagePulls: true
streamingConnectionIdleTimeout: 4h0m0s
syncFrequency: 1m0s
volumeStatsAggPeriod: 1m0s
`

const kubeadmYaml = `apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
clusterName: "{{.ClusterName}}"
kubernetesVersion: "{{.ClusterVersion}}"
controlPlaneEndpoint: "{{.ClusterLbDomain}}:{{.ClusterLbPort}}"
networking:
  dnsDomain: {{.ClusterDnsDomain}}
  serviceSubnet: {{.ClusterSvcCIDR}}
  podSubnet: {{.ClusterPodCIDR}}
apiServer:
  certSANs:
  - 127.0.0.1
  - localhost
  - {{.ClusterApiIP}}
  {{- if ne .ClusterLbIP "127.0.0.1"}}
  - {{.ClusterLbIP}}
  {{- end}}
  - {{.ClusterLbDomain}}
  {{- range .InputCertSANs}}
  - {{.}}
  {{- end}}
  extraArgs:
    feature-gates: TTLAfterFinished=true
    service-node-port-range: {{.ClusterPortRange}}
  extraVolumes:
  - name: localtime
    hostPath: /etc/localtime
    mountPath: /etc/localtime
    readOnly: true
controllerManager:
  extraArgs:
    feature-gates: TTLAfterFinished=true
    experimental-cluster-signing-duration: 876000h
  extraVolumes:
  - hostPath: /etc/localtime
    mountPath: /etc/localtime
    name: localtime
    readOnly: true
scheduler:
  extraArgs:
    feature-gates: TTLAfterFinished=true
  extraVolumes:
  - hostPath: /etc/localtime
    mountPath: /etc/localtime
    name: localtime
    readOnly: true
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: "{{.KubeProxyMode}}"
ipvs:
  minSyncPeriod: 0s
  scheduler: "{{.KubeIPVSScheduler}}"
  syncPeriod: 15s
  {{- if .IsExternalLB}}
  excludeCIDRs: 
  - "{{.ClusterLbIP}}/32"
  {{- end}}
iptables:
  masqueradeAll: true
  masqueradeBit: 14
  minSyncPeriod: 0s
  syncPeriod: 30s
---
`

const kubeadmInitYaml = `apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
bootstrapTokens:
- token: "{{.Token}}"
  description: "kubeadm bootstrap token"
nodeRegistration:
  name: "{{.NodeName}}"
localAPIEndpoint:
  advertiseAddress: "{{.AdvertiseAddress}}"
  bindPort: {{.BindPort}}
certificateKey: {{.CertificateKey}}
---
`

const kubeadmJoinYaml = `apiVersion: kubeadm.k8s.io/v1beta2
kind: JoinConfiguration
discovery:
  bootstrapToken:
    token: {{.Token}}
    apiServerEndpoint: "{{.ClusterLbDomain}}:{{.ClusterLbPort}}"
    caCertHashes: 
    - {{.CaCertHash}}
  timeout: 5m0s
nodeRegistration:
  name: "{{.NodeName}}"
{{- if .IsControlPlane}}
controlPlane:
  localAPIEndpoint:
    advertiseAddress: {{.AdvertiseAddress}}
    bindPort: {{.BindPort}}
  certificateKey: {{.CertificateKey}}
{{- end}}
caCertPath: /etc/kubernetes/pki/ca.crt
---
`

const healthzReaderYaml = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: healthz-reader
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: healthz-reader
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:authenticated
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:unauthenticated
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: healthz-reader
rules:
- nonResourceURLs: ["/healthz", "/healthz/*"] # '*' in a nonResourceURL is a suffix glob match
  verbs: ["get", "post"]
`
