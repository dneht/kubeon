/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package configset

const kubeletYaml = `---
apiVersion: kubelet.config.k8s.io/{{.APIVersion}}
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
containerRuntimeEndpoint: {{.ContainerRuntimeEndpoint}}
enforceNodeAllocatable:
  - pods
allowedUnsafeSysctls:
  - kernel.sem
  - kernel.shm*
  - kernel.msg*
  - fs.mqueue.*
  - net.*
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

const kubeadmYaml = `---
apiVersion: kubeadm.k8s.io/{{.APIVersion}}
kind: ClusterConfiguration
imageRepository: "{{.ImageRepository}}"
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
    {{- if .ClusterFeatureGates}}
    {{.ClusterFeatureGates}}
    {{- end}}
    service-node-port-range: "{{.ClusterPortRange}}"
  extraVolumes:
  - name: localtime
    hostPath: /etc/localtime
    mountPath: /etc/localtime
    readOnly: true
controllerManager:
  extraArgs:
    {{- if .ClusterEnableDual}}
    {{- if gt .ClusterNodeMaskSize 0}}
    node-cidr-mask-size-ipv4: "{{.ClusterNodeMaskSize}}"
    {{- end}}
    {{- if gt .ClusterNodeMaskSizeV6 0}}
    node-cidr-mask-size-ipv6: "{{.ClusterNodeMaskSizeV6}}"
    {{- end}}
    {{- else}}
    {{- if gt .ClusterNodeMaskSize 0}}
    node-cidr-mask-size: "{{.ClusterNodeMaskSize}}"
    {{- end}}
    {{- end}}
    {{- if .ClusterFeatureGates}}
    {{.ClusterFeatureGates}}
    {{- end}}
    {{- if .ClusterSigningDuration}}
    {{.ClusterSigningDuration}}
    {{- end}}
  extraVolumes:
  - hostPath: /etc/localtime
    mountPath: /etc/localtime
    name: localtime
    readOnly: true
scheduler:
  extraArgs:
    {{- if .ClusterFeatureGates}}
    {{.ClusterFeatureGates}}
    {{- end}}
  extraVolumes:
  - hostPath: /etc/localtime
    mountPath: /etc/localtime
    name: localtime
    readOnly: true

---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: {{.ProxyMode}}
ipvs:
  scheduler: {{.IPVSScheduler}}
  minSyncPeriod: 0s
  {{- if .StrictARP}}
  strictARP: {{.StrictARP}}
  {{- end}}
  syncPeriod: 15s
  {{- if .IsExternalLB}}
  excludeCIDRs: 
  - {{.ClusterLbIP}}/32
  {{- end}}
iptables:
  masqueradeAll: true
  masqueradeBit: 14
  minSyncPeriod: 0s
  syncPeriod: 30s

`

const kubeadmInitYaml = `---
apiVersion: kubeadm.k8s.io/{{.APIVersion}}
kind: InitConfiguration
bootstrapTokens:
- token: "{{.Token}}"
  description: "kubeadm bootstrap token"
nodeRegistration:
  name: "{{.NodeName}}"
  {{- if ne .APIVersion "v1beta2"}}
  imagePullPolicy: "{{.ImagePullPolicy}}"
  {{- end}}
localAPIEndpoint:
  advertiseAddress: "{{.AdvertiseAddress}}"
  bindPort: {{.BindPort}}
certificateKey: {{.CertificateKey}}

`

const kubeadmJoinYaml = `---
apiVersion: kubeadm.k8s.io/{{.APIVersion}}
kind: JoinConfiguration
discovery:
  bootstrapToken:
    token: "{{.Token}}"
    apiServerEndpoint: "{{.ClusterLbDomain}}:{{.ClusterLbPort}}"
    caCertHashes: 
    - {{.CaCertHash}}
  timeout: 5m0s
nodeRegistration:
  name: "{{.NodeName}}"
  {{- if ne .APIVersion "v1beta2"}}
  imagePullPolicy: "{{.ImagePullPolicy}}"
  {{- end}}
{{- if .IsControlPlane}}
controlPlane:
  localAPIEndpoint:
    advertiseAddress: {{.AdvertiseAddress}}
    bindPort: {{.BindPort}}
  certificateKey: {{.CertificateKey}}
{{- end}}
caCertPath: /etc/kubernetes/pki/ca.crt

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
- nonResourceURLs:
  - /healthz
  - /healthz/*
  verbs:
  - get
  - post

`
