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

const apiserverUpdaterYaml = `apiVersion: v1
kind: Service
metadata:
  name: apiserver
  namespace: kube-system
  labels:
    app.kubernetes.io/part-of: "kubeon"
    app.kubernetes.io/component: "updater"
spec:
  clusterIP: {{.ClusterLbIP}}
  ports:
  - name: https
    port: 6443
    targetPort: 6443
    protocol: TCP

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: apiserver-updater
  labels:
    app.kubernetes.io/part-of: "kubeon"
    app.kubernetes.io/component: "updater"
    rbac.authorization.k8s.io/aggregate-to-view: "true"
rules:
- apiGroups:
    - ""
  resources:
    - "endpoints"
  verbs:
    - "get"
    - "list"
    - "watch"
    - "create"
    - "update"

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: apiserver-updater
  namespace: kube-system
  labels:
    app.kubernetes.io/part-of: "kubeon"
    app.kubernetes.io/component: "updater"
subjects:
  - kind: ServiceAccount
    name: kubeon-apiserver-updater
    namespace: kube-system
    apiGroup: ""
roleRef:
  kind: ClusterRole
  name: apiserver-updater
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubeon-apiserver-updater
  namespace: kube-system
  labels:
    app.kubernetes.io/part-of: "kubeon"
    app.kubernetes.io/component: "updater"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: apiserver-updater
  namespace: kube-system
  labels:
    app.kubernetes.io/part-of: "kubeon"
    app.kubernetes.io/component: "updater"
spec:
  replicas: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: "apiserver-updater"
      app.kubernetes.io/part-of: "kubeon"
      app.kubernetes.io/component: "updater"
  template:
    metadata:
      labels:
        app.kubernetes.io/name: "apiserver-updater"
        app.kubernetes.io/part-of: "kubeon"
        app.kubernetes.io/component: "updater"
    spec:
      serviceAccountName: kubeon-apiserver-updater
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app.kubernetes.io/name
                operator: In
                values:
                - apiserver-updater
              - key: app.kubernetes.io/part-of
                operator: In
                values:
                - kubeon
              - key: app.kubernetes.io/component
                operator: In
                values:
                - updater
            topologyKey: kubernetes.io/hostname
      containers:
      - name: updater
        image: {{.ImageUrl}}
        imagePullPolicy: IfNotPresent
        resources:
          requests:
            memory: 50Mi
            cpu: 50m
          limits:
            memory: 200Mi
            cpu: 200m
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  revisionHistoryLimit: 3
  minReadySeconds: 30
`
