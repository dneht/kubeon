## Install

> sh -c "$(wget https://dl.sre.pub/on/install.sh -q -O -)"

## Support

> kubeon view support

```text
v1.19.4-v1.19.16
v1.20.1-v1.20.13
v1.21.1-v1.21.7
v1.22.1-v1.22.4
```

## Usage
> k8s_ver=v1.22.4

### Vagrant test

> cd test && vagrant up

### Create cluster

```shell
# ssh node0
vagrant ssh node0
# use root with password 4567890123
su - root
# install kubeon
sh -c "$(wget https://dl.sre.pub/on/install.sh -q -O -)"
# create cluster
# cluster name is "test"
kubeon create test ${k8s_ver} \
    -m 172.20.0.21 \
    -m 172.20.0.22 \
    -m 172.20.0.23 \
    --master-name test10 \
    --master-name test20 \
    --master-name test30 \
    -w 172.20.0.25 \
    --worker-name test50 \
    --default-passwd 4567890123 \
    --interface enp0s8 \
    --log-level debug
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

only support 1.21.x or later

### Add node

```shell
# add one master
kubeon addon test \
    -m 172.20.0.24 \
    --master-name test40 \
    --default-passwd 4567890123 \
    --log-level debug
# add one worker
kubeon addon test \
    -w 172.20.0.26 \
    --worker-name test110 \
    --default-passwd 4567890123 \
    --log-level debug
```

### Del node

```shell
# delon one node
kubeon delon test \
    ip=172.20.0.24 \
    --log-level debug
# or use hostname
kubeon delon test \
    name=test40 \
    --log-level debug
```

### Upgrade cluster

```shell
kubeon upgrade test ${k8s_ver} \
    --log-level debug
```

### Destroy cluster

```shell
kubeon destroy test \
    --log-level debug
```

### Cluster info

```shell
# cluster info
kubeon view info test
# all node ipvs rule
kubeon exec test@all "ipvsadm -ln" -R
```
