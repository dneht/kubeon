## Install

> sh -c "$(wget https://back.pub/kubeon/install.sh -q -O -)"

## Support

> kubeon view support

```text
v1.19.16
v1.20.15
v1.21.14
v1.22.16-v1.22.17
v1.23.17
v1.24.17
v1.25.14-v1.25.16
v1.26.14
v1.27.9-v1.27.11
v1.28.5-v1.28.7
v1.29.1-v1.29.2
```

## Component

> kubeon view cp v1.29.2

```json
{
  "kubernetes": "v1.29.2",
  "pause": "3.9",
  "etcd": "3.5.10",
  "coredns": "1.11.1",
  "crictl": "v1.29.0",
  "runc": "v1.1.12",
  "containerd": "1.7.13",
  "docker": "25.0.3",
  "nvidia": "v3.14.0",
  "kata": "3.2.0",
  "cni": "v1.4.0",
  "calico": "v3.26.4",
  "cilium": "v1.14.7",
  "hubble": "v0.13.0",
  "contour": "v1.28.1",
  "istio": "1.19.7",
  "haproxy": "2.8.6",
  "kruise": "v1.5.2",
  "offline": "20240222",
  "tools": "20240222"
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

online mode is the default and uses `uhub.service.ucloud.cn` as the default mirror source, you can set `--mirror=no` to use `registry.k8s.io` source

offline mode(**--offline**) will download all images on the central machine and import them on each machine

you can try setting the **--mirror** parameter like:

- yes or true: use `uhub.service.ucloud.cn`, **default**
- no or false: use `registry.k8s.io`, if you can access directly
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
kubeon create test v1.28.1 \
    --cilium-enable-dsr \
    -m 192.168.60.21 \
    -m 192.168.60.22 \
    -m 192.168.60.23 \
    --master-name test10 \
    --master-name test20 \
    --master-name test30 \
    -w 192.168.60.24 \
    --worker-name test40 \
    --default-passwd 123456 \
    --interface enp0s8 \
    --ic contour \
    --with-kata \
    --with-kruise \
    --offline \
    --v 6
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
W0915 02:54:55.722472   28148 configset.go:78] Warning: No kubeproxy.config.k8s.io/v1alpha1 config is loaded. Continuing without it: configmaps "kube-proxy" not found

CERTIFICATE                EXPIRES                  RESIDUAL TIME   CERTIFICATE AUTHORITY   EXTERNALLY MANAGED
admin.conf                 Aug 22, 2123 02:36 UTC   99y             ca                      no      
apiserver                  Aug 22, 2123 02:36 UTC   99y             ca                      no      
apiserver-etcd-client      Aug 22, 2123 02:36 UTC   99y             etcd-ca                 no      
apiserver-kubelet-client   Aug 22, 2123 02:36 UTC   99y             ca                      no      
controller-manager.conf    Aug 22, 2123 02:36 UTC   99y             ca                      no      
etcd-healthcheck-client    Aug 22, 2123 02:36 UTC   99y             etcd-ca                 no      
etcd-peer                  Aug 22, 2123 02:36 UTC   99y             etcd-ca                 no      
etcd-server                Aug 22, 2123 02:36 UTC   99y             etcd-ca                 no      
front-proxy-client         Aug 22, 2123 02:36 UTC   99y             front-proxy-ca          no      
scheduler.conf             Aug 22, 2123 02:36 UTC   99y             ca                      no      

CERTIFICATE AUTHORITY   EXPIRES                  RESIDUAL TIME   EXTERNALLY MANAGED
ca                      Aug 22, 2123 02:36 UTC   99y             no      
etcd-ca                 Aug 22, 2123 02:36 UTC   99y             no      
front-proxy-ca          Aug 22, 2123 02:36 UTC   99y             no      
====================node info====================
NAME     STATUS   ROLES           AGE   VERSION   INTERNAL-IP     EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
test10   Ready    control-plane   18m   v1.28.1   192.168.60.21   <none>        Ubuntu 20.04.1 LTS   5.4.0-162-generic   containerd://1.7.5
test20   Ready    control-plane   17m   v1.28.1   192.168.60.22   <none>        Ubuntu 20.04.1 LTS   5.4.0-162-generic   containerd://1.7.5
test30   Ready    control-plane   16m   v1.28.1   192.168.60.23   <none>        Ubuntu 20.04.1 LTS   5.4.0-162-generic   containerd://1.7.5
test40   Ready    worker          14m   v1.28.1   192.168.60.24   <none>        Ubuntu 20.04.1 LTS   5.4.0-162-generic   containerd://1.7.5
====================pod info====================
NAMESPACE        NAME                                         READY   STATUS      RESTARTS      AGE   IP              NODE     NOMINATED NODE   READINESS GATES
kata-system      kata-deploy-km89x                            1/1     Running     0             14m   10.96.2.139     test30   <none>           <none>
kata-system      kata-deploy-n42js                            1/1     Running     0             14m   10.96.1.176     test20   <none>           <none>
kata-system      kata-deploy-pql2r                            1/1     Running     0             14m   10.96.3.82      test40   <none>           <none>
kata-system      kata-deploy-qjl9p                            1/1     Running     0             14m   10.96.0.144     test10   <none>           <none>
kruise-system    kruise-controller-manager-69874777cc-2wwcm   1/1     Running     1             14m   10.96.2.188     test30   <none>           <none>
kruise-system    kruise-controller-manager-69874777cc-jsl6w   1/1     Running     0             14m   10.96.3.215     test40   <none>           <none>
kruise-system    kruise-daemon-5psrs                          1/1     Running     0             14m   192.168.60.22   test20   <none>           <none>
kruise-system    kruise-daemon-c4z2t                          1/1     Running     0             14m   192.168.60.24   test40   <none>           <none>
kruise-system    kruise-daemon-dc6zw                          1/1     Running     1 (13m ago)   14m   192.168.60.21   test10   <none>           <none>
kruise-system    kruise-daemon-f6c5l                          1/1     Running     1             14m   192.168.60.23   test30   <none>           <none>
kube-system      cilium-572kw                                 1/1     Running     0             16m   192.168.60.23   test30   <none>           <none>
kube-system      cilium-c8vlj                                 1/1     Running     0             17m   192.168.60.21   test10   <none>           <none>
kube-system      cilium-df6cj                                 1/1     Running     0             14m   192.168.60.24   test40   <none>           <none>
kube-system      cilium-operator-7b9f4d96fb-6qdmx             1/1     Running     1 (17m ago)   18m   192.168.60.21   test10   <none>           <none>
kube-system      cilium-z5jb6                                 1/1     Running     0             17m   192.168.60.22   test20   <none>           <none>
kube-system      coredns-5dd5756b68-pp8gx                     1/1     Running     0             18m   10.96.0.198     test10   <none>           <none>
kube-system      coredns-5dd5756b68-ppvk5                     1/1     Running     0             18m   10.96.0.107     test10   <none>           <none>
kube-system      etcd-test10                                  1/1     Running     5             18m   192.168.60.21   test10   <none>           <none>
kube-system      etcd-test20                                  1/1     Running     1             17m   192.168.60.22   test20   <none>           <none>
kube-system      etcd-test30                                  1/1     Running     0             16m   192.168.60.23   test30   <none>           <none>
kube-system      hubble-relay-86894b7788-99czd                1/1     Running     0             17m   10.96.0.214     test10   <none>           <none>
kube-system      hubble-ui-7d6f59fbf5-tm8z8                   2/2     Running     0             17m   10.96.0.220     test10   <none>           <none>
kube-system      kube-apiserver-test10                        1/1     Running     5             18m   192.168.60.21   test10   <none>           <none>
kube-system      kube-apiserver-test20                        1/1     Running     3             17m   192.168.60.22   test20   <none>           <none>
kube-system      kube-apiserver-test30                        1/1     Running     4 (16m ago)   16m   192.168.60.23   test30   <none>           <none>
kube-system      kube-controller-manager-test10               1/1     Running     9 (17m ago)   18m   192.168.60.21   test10   <none>           <none>
kube-system      kube-controller-manager-test20               1/1     Running     3             17m   192.168.60.22   test20   <none>           <none>
kube-system      kube-controller-manager-test30               1/1     Running     3             15m   192.168.60.23   test30   <none>           <none>
kube-system      kube-scheduler-test10                        1/1     Running     9 (17m ago)   18m   192.168.60.21   test10   <none>           <none>
kube-system      kube-scheduler-test20                        1/1     Running     3             17m   192.168.60.22   test20   <none>           <none>
kube-system      kube-scheduler-test30                        1/1     Running     3             16m   192.168.60.23   test30   <none>           <none>
kube-system      local-haproxy-test40                         1/1     Running     0             14m   192.168.60.24   test40   <none>           <none>
projectcontour   contour-574bcf6b5d-8thns                     1/1     Running     0             14m   10.96.3.221     test40   <none>           <none>
projectcontour   contour-574bcf6b5d-wrt69                     1/1     Running     0             14m   10.96.1.235     test20   <none>           <none>
projectcontour   contour-certgen-v1.25.2-2pvhn                0/1     Completed   0             14m   10.96.3.89      test40   <none>           <none>
projectcontour   envoy-2777q                                  2/2     Running     0             14m   10.96.2.191     test30   <none>           <none>
projectcontour   envoy-kzzx7                                  2/2     Running     0             14m   10.96.3.152     test40   <none>           <none>
projectcontour   envoy-vcwxp                                  2/2     Running     0             14m   10.96.1.24      test20   <none>           <none>
projectcontour   envoy-vn4k4                                  2/2     Running     0             14m   10.96.0.64      test10   <none>           <none>
====================container info====================

kata-deploy-km89x:      quay.io/kata-containers/kata-deploy:3.1.3, 
kata-deploy-n42js:      quay.io/kata-containers/kata-deploy:3.1.3, 
kata-deploy-pql2r:      quay.io/kata-containers/kata-deploy:3.1.3, 
kata-deploy-qjl9p:      quay.io/kata-containers/kata-deploy:3.1.3, 
kruise-controller-manager-69874777cc-2wwcm:     openkruise/kruise-manager:v1.5.0, 
kruise-controller-manager-69874777cc-jsl6w:     openkruise/kruise-manager:v1.5.0, 
kruise-daemon-5psrs:    openkruise/kruise-manager:v1.5.0, 
kruise-daemon-c4z2t:    openkruise/kruise-manager:v1.5.0, 
kruise-daemon-dc6zw:    openkruise/kruise-manager:v1.5.0, 
kruise-daemon-f6c5l:    openkruise/kruise-manager:v1.5.0, 
cilium-572kw:   quay.io/cilium/cilium:v1.13.6, 
cilium-c8vlj:   quay.io/cilium/cilium:v1.13.6, 
cilium-df6cj:   quay.io/cilium/cilium:v1.13.6, 
cilium-operator-7b9f4d96fb-6qdmx:       quay.io/cilium/operator-generic:v1.13.6, 
cilium-z5jb6:   quay.io/cilium/cilium:v1.13.6, 
coredns-5dd5756b68-pp8gx:       registry.k8s.io/coredns/coredns:v1.10.1, 
coredns-5dd5756b68-ppvk5:       registry.k8s.io/coredns/coredns:v1.10.1, 
etcd-test10:    registry.k8s.io/etcd:3.5.9-0, 
etcd-test20:    registry.k8s.io/etcd:3.5.9-0, 
etcd-test30:    registry.k8s.io/etcd:3.5.9-0, 
hubble-relay-86894b7788-99czd:  quay.io/cilium/hubble-relay:v1.13.6, 
hubble-ui-7d6f59fbf5-tm8z8:     quay.io/cilium/hubble-ui:v0.12.0, quay.io/cilium/hubble-ui-backend:v0.12.0, 
kube-apiserver-test10:  registry.k8s.io/kube-apiserver:v1.28.1, 
kube-apiserver-test20:  registry.k8s.io/kube-apiserver:v1.28.1, 
kube-apiserver-test30:  registry.k8s.io/kube-apiserver:v1.28.1, 
kube-controller-manager-test10: registry.k8s.io/kube-controller-manager:v1.28.1, 
kube-controller-manager-test20: registry.k8s.io/kube-controller-manager:v1.28.1, 
kube-controller-manager-test30: registry.k8s.io/kube-controller-manager:v1.28.1, 
kube-scheduler-test10:  registry.k8s.io/kube-scheduler:v1.28.1, 
kube-scheduler-test20:  registry.k8s.io/kube-scheduler:v1.28.1, 
kube-scheduler-test30:  registry.k8s.io/kube-scheduler:v1.28.1, 
local-haproxy-test40:   kubeon/local-haproxy:v1.28.1, 
contour-574bcf6b5d-8thns:       ghcr.io/projectcontour/contour:v1.25.2, 
contour-574bcf6b5d-wrt69:       ghcr.io/projectcontour/contour:v1.25.2, 
contour-certgen-v1.25.2-2pvhn:  ghcr.io/projectcontour/contour:v1.25.2, 
envoy-2777q:    ghcr.io/projectcontour/contour:v1.25.2, envoyproxy/envoy:v1.26.4, 
envoy-kzzx7:    ghcr.io/projectcontour/contour:v1.25.2, envoyproxy/envoy:v1.26.4, 
envoy-vcwxp:    ghcr.io/projectcontour/contour:v1.25.2, envoyproxy/envoy:v1.26.4, 
envoy-vn4k4:    ghcr.io/projectcontour/contour:v1.25.2, envoyproxy/envoy:v1.26.4, 
====================etcd info====================
using etcdctl version: 3.5.9
+------------------+---------+--------+----------------------------+----------------------------+------------+
|        ID        | STATUS  |  NAME  |         PEER ADDRS         |        CLIENT ADDRS        | IS LEARNER |
+------------------+---------+--------+----------------------------+----------------------------+------------+
| 9612592669ab97fd | started | test20 | https://192.168.60.22:2380 | https://192.168.60.22:2379 |      false |
| af3f764e67c48c30 | started | test10 | https://192.168.60.21:2380 | https://192.168.60.21:2379 |      false |
| dae5dde33f2c1ab3 | started | test30 | https://192.168.60.23:2380 | https://192.168.60.23:2379 |      false |
+------------------+---------+--------+----------------------------+----------------------------+------------+
```