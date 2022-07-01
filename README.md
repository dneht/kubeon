## Install

> sh -c "$(wget https://back.pub/kubeon/install.sh -q -O -)"

## Support

> kubeon view support

```text
v1.19.4-v1.19.16
v1.20.1-v1.20.15
v1.21.1-v1.21.14
v1.22.1-v1.22.11
v1.23.1-v1.23.8
v1.24.1-v1.24.2
```
## Component

> kubeon view component v1.24.2

```json
{
  "kubernetes": "v1.24.2",
  "etcd": "3.5.3",
  "coredns": "1.9.3",
  "crictl": "v1.24.2",
  "runc": "v1.1.3",
  "containerd": "1.6.6",
  "docker": "20.10.17",
  "nvidia": "v3.10.0",
  "kata": "2.4.2",
  "cni": "v1.1.1",
  "calico": "v3.22.3",
  "cilium": "v1.11.6",
  "contour": "v1.21.1",
  "haproxy": "2.5.7"
}
```

## Modified

you can use the `kubeon view component` to get detailed information

### kubeadm

kubeadm was recompiled and only the certificate generation part was modified

### Kubelet

kubelet will add the following parameters and use systemd

1. allowed-unsafe-sysctls=kernel.sem,kernel.shm*,kernel.msg*,fs.mqueue.*,net.*

### Coredns

coredns is always the latest version, example

- 1.23.8 actually uses version 1.8.7
- 1.24.2 actually uses version 1.9.3

other images such as etcd remain the same

## Offline

online mode is the default and uses `registry.cn-hangzhou.aliyuncs.com` as the default mirror source, you can set `--mirror=no` to use `k8s.gcr.io` source

offline mode(**--offline**) will download all images on the central machine and import them on each machine

you can try setting the **--mirror** parameter like:

- yes or true: use `registry.cn-hangzhou.aliyuncs.com`, **default**
- no or false: use `k8s.gcr.io`, if you can access directly
- any other docker mirror address

## Usage

### Vagrant test

> cd test && vagrant up

### Create cluster

```shell
# ssh node0
vagrant ssh node0
# use root with password 123456
su - root
# install kubeon
sh -c "$(wget https://dl.sre.pub/on/install.sh -q -O -)"
# create cluster
# cluster name is "test"
kubeon create test ${k8s_ver} \
  -m 192.168.60.21 \
  -m 192.168.60.22 \
  -m 192.168.60.23 \
  --master-name test10 \
  --master-name test20 \
  --master-name test30 \
  -w 192.168.60.25 \
  --worker-name test50 \
  --default-passwd 123456 \
  --ic contour \
  --interface enp0s8 \
  --v 4
```
Please use `source /etc/profile` for autocompletion

#### --cri
default is `containerd`, you can also use `docker`

#### --cni
default is `calico`, you can also use `none` and install cni later

#### --lb-ip && --lb-port
external load balancer ip and port

#### --lb-mode 
if you set lb-ip, the following two inner ha schemes are invalid

##### haproxy
default is `haproxy`, will be created a static pod `local-haproxy` on each worker 

##### updater
create one deployment on now cluster, and one service on each worker

#### --ic
default ingress is `none`, you can also use `contour`

only support v1.21.1 or later

#### --with-nvidia
default is true, if `test -c /dev/nvidia*` is true, `nvidia-container-runtime` will be installed instead of `runc`

only support (debian9+, ubuntu16+) and containerd and (v1.22.1 or later)

#### --with-kata
default is false, if true kata will be installed in namespace `kata-system` using `kata-deploy`

only support v1.22.1 or later

### Add node

```shell
# add one master
kubeon add test \
  -m 192.168.60.26 \
  --master-name test60 \
  --default-passwd 123456 \
  --v 4
# add one worker
kubeon add test \
  -w 192.168.60.24 \
  --worker-name test40 \
  --default-passwd 123456 \
  --v 4
```

### Del node

```shell
# delon one node
kubeon del test ip=192.168.60.24 --v 4
# or use hostname
kubeon del test name=test40,test50 --v 4
```

### Upgrade cluster

```shell
kubeon upgrade test ${k8s_ver} --v 4
```

### Destroy cluster

```shell
kubeon destroy test --v 4
```

### Cluster info

```shell
# etcd member list
kubeon etcd test member list
# cluster info
kubeon display test
# all node ipvs rule
kubeon exec test@all "ipvsadm -ln" -R
```

#### example

##### creat new cluster

```bash
kubeon create test v1.22.6 \
-m 192.168.60.21 \
-m 192.168.60.22 \
-m 192.168.60.23 \
--master-name test10 \
--master-name test20 \
--master-name test30 \
-w 192.168.60.24 \
--worker-name test40 \
--default-passwd 123456 \
--ic contour \
--with-kata \
--interface enp0s8 \
--v 4
```

Note:

> --offline(bool) default is false

> --mirror(string) default is yes, will use alibaba cloud mirror

##### display cluster info

```bash
kubeon display test 
```

```text
====================cert info====================
[check-expiration] Reading configuration from the cluster...
[check-expiration] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -o yaml'

CERTIFICATE                EXPIRES                  RESIDUAL TIME   CERTIFICATE AUTHORITY   EXTERNALLY MANAGED
admin.conf                 Jan 09, 2122 13:32 UTC   99y             ca                      no      
apiserver                  Jan 09, 2122 13:32 UTC   99y             ca                      no      
apiserver-etcd-client      Jan 09, 2122 13:32 UTC   99y             etcd-ca                 no      
apiserver-kubelet-client   Jan 09, 2122 13:32 UTC   99y             ca                      no      
controller-manager.conf    Jan 09, 2122 13:32 UTC   99y             ca                      no      
etcd-healthcheck-client    Jan 09, 2122 13:32 UTC   99y             etcd-ca                 no      
etcd-peer                  Jan 09, 2122 13:32 UTC   99y             etcd-ca                 no      
etcd-server                Jan 09, 2122 13:32 UTC   99y             etcd-ca                 no      
front-proxy-client         Jan 09, 2122 13:32 UTC   99y             front-proxy-ca          no      
scheduler.conf             Jan 09, 2122 13:32 UTC   99y             ca                      no      

CERTIFICATE AUTHORITY   EXPIRES                  RESIDUAL TIME   EXTERNALLY MANAGED
ca                      Jan 09, 2122 13:32 UTC   99y             no      
etcd-ca                 Jan 09, 2122 13:32 UTC   99y             no      
front-proxy-ca          Jan 09, 2122 13:32 UTC   99y             no      
====================node info====================
NAME     STATUS   ROLES                  AGE    VERSION   INTERNAL-IP     EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
test10   Ready    control-plane,master   121m   v1.22.6   192.168.60.21   <none>        Ubuntu 20.04.1 LTS   5.4.0-54-generic   containerd://1.5.9
test20   Ready    control-plane,master   120m   v1.22.6   192.168.60.22   <none>        Ubuntu 20.04.1 LTS   5.4.0-54-generic   containerd://1.5.9
test30   Ready    control-plane,master   118m   v1.22.6   192.168.60.23   <none>        Ubuntu 20.04.1 LTS   5.4.0-54-generic   containerd://1.5.9
test40   Ready    worker                 118m   v1.22.6   192.168.60.24   <none>        Ubuntu 20.04.1 LTS   5.4.0-54-generic   containerd://1.5.9
====================pod info====================
NAMESPACE        NAME                                       READY   STATUS    RESTARTS        AGE     IP              NODE     NOMINATED NODE   READINESS GATES
kata-system      kata-deploy-2mbdt                          1/1     Running   0               107m    10.107.66.3     test30   <none>           <none>
kata-system      kata-deploy-7tg6c                          1/1     Running   0               107m    10.108.198.72   test40   <none>           <none>
kata-system      kata-deploy-bbnp7                          1/1     Running   0               107m    10.96.207.9     test10   <none>           <none>
kata-system      kata-deploy-fsjvv                          1/1     Running   0               107m    10.111.65.197   test20   <none>           <none>
kube-system      calico-kube-controllers-5d54f88696-th8mn   1/1     Running   0               3m19s   192.168.60.21   test10   <none>           <none>
kube-system      calico-node-9qczk                          1/1     Running   0               118m    192.168.60.23   test30   <none>           <none>
kube-system      calico-node-bkr7t                          1/1     Running   0               118m    192.168.60.24   test40   <none>           <none>
kube-system      calico-node-lv2xc                          1/1     Running   42 (106s ago)   120m    192.168.60.21   test10   <none>           <none>
kube-system      calico-node-rmhlz                          1/1     Running   0               120m    192.168.60.22   test20   <none>           <none>
kube-system      coredns-7f676b947b-ccf96                   1/1     Running   0               3m19s   10.107.66.7     test30   <none>           <none>
kube-system      coredns-7f676b947b-h9bh2                   1/1     Running   0               6m38s   10.96.207.11    test10   <none>           <none>
kube-system      etcd-test10                                1/1     Running   0               120m    192.168.60.21   test10   <none>           <none>
kube-system      etcd-test20                                1/1     Running   0               108m    192.168.60.22   test20   <none>           <none>
kube-system      etcd-test30                                1/1     Running   0               107m    192.168.60.23   test30   <none>           <none>
kube-system      kube-apiserver-test10                      1/1     Running   0               7m16s   192.168.60.21   test10   <none>           <none>
kube-system      kube-apiserver-test20                      1/1     Running   0               5m27s   192.168.60.22   test20   <none>           <none>
kube-system      kube-apiserver-test30                      1/1     Running   0               3m54s   192.168.60.23   test30   <none>           <none>
kube-system      kube-controller-manager-test10             1/1     Running   0               7m      192.168.60.21   test10   <none>           <none>
kube-system      kube-controller-manager-test20             1/1     Running   0               5m10s   192.168.60.22   test20   <none>           <none>
kube-system      kube-controller-manager-test30             1/1     Running   0               3m41s   192.168.60.23   test30   <none>           <none>
kube-system      kube-proxy-9szkz                           1/1     Running   0               6m42s   192.168.60.24   test40   <none>           <none>
kube-system      kube-proxy-dqvdf                           1/1     Running   0               6m35s   192.168.60.22   test20   <none>           <none>
kube-system      kube-proxy-lnxrx                           1/1     Running   0               6m38s   192.168.60.23   test30   <none>           <none>
kube-system      kube-proxy-m5kct                           1/1     Running   0               6m23s   192.168.60.21   test10   <none>           <none>
kube-system      kube-scheduler-test10                      1/1     Running   0               6m45s   192.168.60.21   test10   <none>           <none>
kube-system      kube-scheduler-test20                      1/1     Running   0               4m55s   192.168.60.22   test20   <none>           <none>
kube-system      kube-scheduler-test30                      1/1     Running   0               3m25s   192.168.60.23   test30   <none>           <none>
kube-system      local-haproxy-test40                       1/1     Running   0               2m49s   192.168.60.24   test40   <none>           <none>
projectcontour   contour-7fffbd4448-qbgcz                   1/1     Running   0               3m19s   10.107.66.6     test30   <none>           <none>
projectcontour   contour-7fffbd4448-qkfzp                   1/1     Running   0               4m58s   10.111.65.199   test20   <none>           <none>
projectcontour   envoy-4bkxn                                2/2     Running   0               104m    10.111.65.198   test20   <none>           <none>
projectcontour   envoy-bd9hs                                2/2     Running   1               106m    10.108.198.73   test40   <none>           <none>
projectcontour   envoy-jp8m9                                2/2     Running   0               106m    10.96.207.10    test10   <none>           <none>
projectcontour   envoy-nqxxd                                2/2     Running   0               103m    10.107.66.5     test30   <none>           <none>
====================container info====================

kata-deploy-2mbdt:      registry.cn-hangzhou.aliyuncs.com/kubeon/kata-deploy:2.3.0, 
kata-deploy-7tg6c:      registry.cn-hangzhou.aliyuncs.com/kubeon/kata-deploy:2.3.0, 
kata-deploy-bbnp7:      registry.cn-hangzhou.aliyuncs.com/kubeon/kata-deploy:2.3.0, 
kata-deploy-fsjvv:      registry.cn-hangzhou.aliyuncs.com/kubeon/kata-deploy:2.3.0, 
calico-kube-controllers-5d54f88696-th8mn:       registry.cn-hangzhou.aliyuncs.com/kubeon/calico-kube-controllers:v3.21.4, 
calico-node-9qczk:      registry.cn-hangzhou.aliyuncs.com/kubeon/calico-node:v3.21.4, 
calico-node-bkr7t:      registry.cn-hangzhou.aliyuncs.com/kubeon/calico-node:v3.21.4, 
calico-node-lv2xc:      registry.cn-hangzhou.aliyuncs.com/kubeon/calico-node:v3.21.4, 
calico-node-rmhlz:      registry.cn-hangzhou.aliyuncs.com/kubeon/calico-node:v3.21.4, 
coredns-7f676b947b-ccf96:       registry.cn-hangzhou.aliyuncs.com/kubeon/coredns:v1.8.4, 
coredns-7f676b947b-h9bh2:       registry.cn-hangzhou.aliyuncs.com/kubeon/coredns:v1.8.4, 
etcd-test10:    registry.cn-hangzhou.aliyuncs.com/kubeon/etcd:3.5.0-0, 
etcd-test20:    registry.cn-hangzhou.aliyuncs.com/kubeon/etcd:3.5.0-0, 
etcd-test30:    registry.cn-hangzhou.aliyuncs.com/kubeon/etcd:3.5.0-0, 
kube-apiserver-test10:  registry.cn-hangzhou.aliyuncs.com/kubeon/kube-apiserver:v1.22.6, 
kube-apiserver-test20:  registry.cn-hangzhou.aliyuncs.com/kubeon/kube-apiserver:v1.22.6, 
kube-apiserver-test30:  registry.cn-hangzhou.aliyuncs.com/kubeon/kube-apiserver:v1.22.6, 
kube-controller-manager-test10: registry.cn-hangzhou.aliyuncs.com/kubeon/kube-controller-manager:v1.22.6, 
kube-controller-manager-test20: registry.cn-hangzhou.aliyuncs.com/kubeon/kube-controller-manager:v1.22.6, 
kube-controller-manager-test30: registry.cn-hangzhou.aliyuncs.com/kubeon/kube-controller-manager:v1.22.6, 
kube-proxy-9szkz:       registry.cn-hangzhou.aliyuncs.com/kubeon/kube-proxy:v1.22.6, 
kube-proxy-dqvdf:       registry.cn-hangzhou.aliyuncs.com/kubeon/kube-proxy:v1.22.6, 
kube-proxy-lnxrx:       registry.cn-hangzhou.aliyuncs.com/kubeon/kube-proxy:v1.22.6, 
kube-proxy-m5kct:       registry.cn-hangzhou.aliyuncs.com/kubeon/kube-proxy:v1.22.6, 
kube-scheduler-test10:  registry.cn-hangzhou.aliyuncs.com/kubeon/kube-scheduler:v1.22.6, 
kube-scheduler-test20:  registry.cn-hangzhou.aliyuncs.com/kubeon/kube-scheduler:v1.22.6, 
kube-scheduler-test30:  registry.cn-hangzhou.aliyuncs.com/kubeon/kube-scheduler:v1.22.6, 
local-haproxy-test40:   registry.cn-hangzhou.aliyuncs.com/kubeon/local-haproxy:v1.22.6, 
contour-7fffbd4448-qbgcz:       registry.cn-hangzhou.aliyuncs.com/kubeon/contour:v1.19.1, 
contour-7fffbd4448-qkfzp:       registry.cn-hangzhou.aliyuncs.com/kubeon/contour:v1.19.1, 
envoy-4bkxn:    registry.cn-hangzhou.aliyuncs.com/kubeon/contour:v1.19.1, registry.cn-hangzhou.aliyuncs.com/kubeon/envoy:v1.19.1, 
envoy-bd9hs:    registry.cn-hangzhou.aliyuncs.com/kubeon/contour:v1.19.1, registry.cn-hangzhou.aliyuncs.com/kubeon/envoy:v1.19.1, 
envoy-jp8m9:    registry.cn-hangzhou.aliyuncs.com/kubeon/contour:v1.19.1, registry.cn-hangzhou.aliyuncs.com/kubeon/envoy:v1.19.1, 
envoy-nqxxd:    registry.cn-hangzhou.aliyuncs.com/kubeon/contour:v1.19.1, registry.cn-hangzhou.aliyuncs.com/kubeon/envoy:v1.19.1, 
====================etcd info====================
using etcdctl version: 3.5.0
+------------------+---------+--------+----------------------------+----------------------------+------------+
|        ID        | STATUS  |  NAME  |         PEER ADDRS         |        CLIENT ADDRS        | IS LEARNER |
+------------------+---------+--------+----------------------------+----------------------------+------------+
| 267f4cbf88a000b1 | started | test20 | https://192.168.60.22:2380 | https://192.168.60.22:2379 |      false |
| af3f764e67c48c30 | started | test10 | https://192.168.60.21:2380 | https://192.168.60.21:2379 |      false |
| d969b5b381f4d03b | started | test30 | https://192.168.60.23:2380 | https://192.168.60.23:2379 |      false |
+------------------+---------+--------+----------------------------+----------------------------+------------+
```