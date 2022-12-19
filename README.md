## Install

> sh -c "$(wget https://back.pub/kubeon/install.sh -q -O -)"

## Support

> kubeon view support

```text
v1.19.4-v1.19.16
v1.20.1-v1.20.15
v1.21.1-v1.21.14
v1.22.1-v1.22.15
v1.23.1-v1.23.12
v1.24.1-v1.24.6
v1.25.1-v1.25.2
```

v1.19.4-v1.19.15

## Component

> kubeon view component v1.25.2

```json
{
  "kubernetes": "v1.25.2",
  "etcd": "3.5.4",
  "coredns": "1.9.4",
  "crictl": "v1.25.0",
  "runc": "v1.1.4",
  "containerd": "1.6.8",
  "docker": "20.10.18",
  "nvidia": "v3.11.0",
  "kata": "2.5.1",
  "cni": "v1.1.1",
  "calico": "v3.23.3",
  "cilium": "v1.12.2",
  "contour": "v1.22.1",
  "haproxy": "2.6.5"
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

- 1.23.x actually uses version 1.8.7
- 1.24.x actually uses version 1.9.3

other images such as etcd remain the same

## Offline

online mode is the default and uses `registry.cn-hangzhou.aliyuncs.com` as the default mirror source, you can set `--mirror=no` to use `k8s.gcr.io` source

offline mode(**--offline**) will download all images on the central machine and import them on each machine

you can try setting the **--mirror** parameter like:

- yes or true: use `registry.cn-hangzhou.aliyuncs.com`, **default**
- no or false: use `k8s.gcr.io`, if you can access directly
- any other docker mirror address, like `mirror.ccs.tencentyun.com` if tencent is used

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

only support v1.22.1 or later with containerd

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
admin.conf                 Sep 03, 2122 06:52 UTC   99y             ca                      no      
apiserver                  Sep 03, 2122 06:52 UTC   99y             ca                      no      
apiserver-etcd-client      Sep 03, 2122 06:52 UTC   99y             etcd-ca                 no      
apiserver-kubelet-client   Sep 03, 2122 06:52 UTC   99y             ca                      no      
controller-manager.conf    Sep 03, 2122 06:52 UTC   99y             ca                      no      
etcd-healthcheck-client    Sep 03, 2122 06:52 UTC   99y             etcd-ca                 no      
etcd-peer                  Sep 03, 2122 06:52 UTC   99y             etcd-ca                 no      
etcd-server                Sep 03, 2122 06:52 UTC   99y             etcd-ca                 no      
front-proxy-client         Sep 03, 2122 06:52 UTC   99y             front-proxy-ca          no      
scheduler.conf             Sep 03, 2122 06:52 UTC   99y             ca                      no      

CERTIFICATE AUTHORITY   EXPIRES                  RESIDUAL TIME   EXTERNALLY MANAGED
ca                      Sep 03, 2122 06:52 UTC   99y             no      
etcd-ca                 Sep 03, 2122 06:52 UTC   99y             no      
front-proxy-ca          Sep 03, 2122 06:52 UTC   99y             no      
====================node info====================
NAME     STATUS   ROLES           AGE    VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE           KERNEL-VERSION      CONTAINER-RUNTIME
test10   Ready    control-plane   3h1m   v1.24.6   172.30.0.94    <none>        Ubuntu 22.04 LTS   5.15.0-40-generic   containerd://1.6.8
test11   Ready    control-plane   3h     v1.24.6   172.30.0.121   <none>        Ubuntu 22.04 LTS   5.15.0-40-generic   containerd://1.6.8
test12   Ready    control-plane   179m   v1.24.6   172.30.0.87    <none>        Ubuntu 22.04 LTS   5.15.0-40-generic   containerd://1.6.8
test13   Ready    worker          178m   v1.24.6   172.30.0.112   <none>        Ubuntu 22.04 LTS   5.15.0-40-generic   containerd://1.6.8
====================pod info====================
NAMESPACE        NAME                                      READY   STATUS    RESTARTS      AGE     IP              NODE     NOMINATED NODE   READINESS GATES
kube-system      calico-kube-controllers-ccdc8c787-rcvwg   1/1     Running   0             112s    172.30.0.87     test12   <none>           <none>
kube-system      calico-node-gdsgm                         1/1     Running   0             63m     172.30.0.112    test13   <none>           <none>
kube-system      calico-node-thnqc                         1/1     Running   0             62m     172.30.0.87     test12   <none>           <none>
kube-system      calico-node-v6qfb                         1/1     Running   0             62m     172.30.0.121    test11   <none>           <none>
kube-system      calico-node-vqnml                         1/1     Running   0             61m     172.30.0.94     test10   <none>           <none>
kube-system      coredns-599ddf989b-jrrvk                  1/1     Running   0             4m51s   10.96.207.16    test10   <none>           <none>
kube-system      coredns-599ddf989b-m7cjt                  1/1     Running   0             112s    10.96.139.12    test12   <none>           <none>
kube-system      etcd-test10                               1/1     Running   2 (68m ago)   67m     172.30.0.94     test10   <none>           <none>
kube-system      etcd-test11                               1/1     Running   1 (65m ago)   65m     172.30.0.121    test11   <none>           <none>
kube-system      etcd-test12                               1/1     Running   1 (63m ago)   63m     172.30.0.87     test12   <none>           <none>
kube-system      kube-apiserver-test10                     1/1     Running   0             5m40s   172.30.0.94     test10   <none>           <none>
kube-system      kube-apiserver-test11                     1/1     Running   0             3m50s   172.30.0.121    test11   <none>           <none>
kube-system      kube-apiserver-test12                     1/1     Running   0             2m24s   172.30.0.87     test12   <none>           <none>
kube-system      kube-controller-manager-test10            1/1     Running   0             5m20s   172.30.0.94     test10   <none>           <none>
kube-system      kube-controller-manager-test11            1/1     Running   0             3m35s   172.30.0.121    test11   <none>           <none>
kube-system      kube-controller-manager-test12            1/1     Running   0             2m11s   172.30.0.87     test12   <none>           <none>
kube-system      kube-proxy-7z6gl                          1/1     Running   0             4m40s   172.30.0.112    test13   <none>           <none>
kube-system      kube-proxy-8tc2w                          1/1     Running   0             4m59s   172.30.0.87     test12   <none>           <none>
kube-system      kube-proxy-f7b96                          1/1     Running   0             4m44s   172.30.0.121    test11   <none>           <none>
kube-system      kube-proxy-lljzg                          1/1     Running   0             4m49s   172.30.0.94     test10   <none>           <none>
kube-system      kube-scheduler-test10                     1/1     Running   0             5m4s    172.30.0.94     test10   <none>           <none>
kube-system      kube-scheduler-test11                     1/1     Running   0             3m23s   172.30.0.121    test11   <none>           <none>
kube-system      kube-scheduler-test12                     1/1     Running   0             116s    172.30.0.87     test12   <none>           <none>
kube-system      local-haproxy-test13                      1/1     Running   0             60s     172.30.0.112    test13   <none>           <none>
projectcontour   contour-57cb7f4d56-k4ghb                  1/1     Running   0             112s    10.96.207.17    test10   <none>           <none>
projectcontour   contour-57cb7f4d56-qkssc                  1/1     Running   0             3m18s   10.107.89.76    test11   <none>           <none>
projectcontour   envoy-9gzz2                               2/2     Running   0             168m    10.96.139.5     test12   <none>           <none>
projectcontour   envoy-9lxsn                               2/2     Running   0             167m    10.96.207.10    test10   <none>           <none>
projectcontour   envoy-g6r27                               2/2     Running   0             168m    10.110.246.71   test13   <none>           <none>
projectcontour   envoy-krklb                               2/2     Running   0             168m    10.107.89.68    test11   <none>           <none>
====================container info====================

calico-kube-controllers-ccdc8c787-rcvwg:	mirror.ccs.tencentyun.com/kubeon/calico-kube-controllers:v3.23.3, 
calico-node-gdsgm:	mirror.ccs.tencentyun.com/kubeon/calico-node:v3.23.3, 
calico-node-thnqc:	mirror.ccs.tencentyun.com/kubeon/calico-node:v3.23.3, 
calico-node-v6qfb:	mirror.ccs.tencentyun.com/kubeon/calico-node:v3.23.3, 
calico-node-vqnml:	mirror.ccs.tencentyun.com/kubeon/calico-node:v3.23.3, 
coredns-599ddf989b-jrrvk:	mirror.ccs.tencentyun.com/kubeon/coredns:v1.8.6, 
coredns-599ddf989b-m7cjt:	mirror.ccs.tencentyun.com/kubeon/coredns:v1.8.6, 
etcd-test10:	mirror.ccs.tencentyun.com/kubeon/etcd:3.5.3-0, 
etcd-test11:	mirror.ccs.tencentyun.com/kubeon/etcd:3.5.3-0, 
etcd-test12:	mirror.ccs.tencentyun.com/kubeon/etcd:3.5.3-0, 
kube-apiserver-test10:	mirror.ccs.tencentyun.com/kubeon/kube-apiserver:v1.24.6, 
kube-apiserver-test11:	mirror.ccs.tencentyun.com/kubeon/kube-apiserver:v1.24.6, 
kube-apiserver-test12:	mirror.ccs.tencentyun.com/kubeon/kube-apiserver:v1.24.6, 
kube-controller-manager-test10:	mirror.ccs.tencentyun.com/kubeon/kube-controller-manager:v1.24.6, 
kube-controller-manager-test11:	mirror.ccs.tencentyun.com/kubeon/kube-controller-manager:v1.24.6, 
kube-controller-manager-test12:	mirror.ccs.tencentyun.com/kubeon/kube-controller-manager:v1.24.6, 
kube-proxy-7z6gl:	mirror.ccs.tencentyun.com/kubeon/kube-proxy:v1.24.6, 
kube-proxy-8tc2w:	mirror.ccs.tencentyun.com/kubeon/kube-proxy:v1.24.6, 
kube-proxy-f7b96:	mirror.ccs.tencentyun.com/kubeon/kube-proxy:v1.24.6, 
kube-proxy-lljzg:	mirror.ccs.tencentyun.com/kubeon/kube-proxy:v1.24.6, 
kube-scheduler-test10:	mirror.ccs.tencentyun.com/kubeon/kube-scheduler:v1.24.6, 
kube-scheduler-test11:	mirror.ccs.tencentyun.com/kubeon/kube-scheduler:v1.24.6, 
kube-scheduler-test12:	mirror.ccs.tencentyun.com/kubeon/kube-scheduler:v1.24.6, 
local-haproxy-test13:	mirror.ccs.tencentyun.com/kubeon/local-haproxy:v1.24.6, 
contour-57cb7f4d56-k4ghb:	mirror.ccs.tencentyun.com/kubeon/contour:v1.22.1, 
contour-57cb7f4d56-qkssc:	mirror.ccs.tencentyun.com/kubeon/contour:v1.22.1, 
envoy-9gzz2:	mirror.ccs.tencentyun.com/kubeon/contour:v1.22.1, mirror.ccs.tencentyun.com/kubeon/envoy:v1.23.1, 
envoy-9lxsn:	mirror.ccs.tencentyun.com/kubeon/contour:v1.22.1, mirror.ccs.tencentyun.com/kubeon/envoy:v1.23.1, 
envoy-g6r27:	mirror.ccs.tencentyun.com/kubeon/contour:v1.22.1, mirror.ccs.tencentyun.com/kubeon/envoy:v1.23.1, 
envoy-krklb:	mirror.ccs.tencentyun.com/kubeon/contour:v1.22.1, mirror.ccs.tencentyun.com/kubeon/envoy:v1.23.1, 
====================etcd info====================
using etcdctl version: 3.5.3
+------------------+---------+--------+---------------------------+---------------------------+------------+
|        ID        | STATUS  |  NAME  |        PEER ADDRS         |       CLIENT ADDRS        | IS LEARNER |
+------------------+---------+--------+---------------------------+---------------------------+------------+
|  4290151dfe289d7 | started | test10 |  https://172.30.0.94:2380 |  https://172.30.0.94:2379 |      false |
|  f6b5a58e0c5d32a | started | test11 | https://172.30.0.121:2380 | https://172.30.0.121:2379 |      false |
| f51f3126cf6f26c4 | started | test12 |  https://172.30.0.87:2380 |  https://172.30.0.87:2379 |      false |
+------------------+---------+--------+---------------------------+---------------------------+------------+
```