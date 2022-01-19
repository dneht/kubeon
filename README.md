## Install

> sh -c "$(wget https://dl.sre.pub/on/install.sh -q -O -)"

## Support

> kubeon view support

```text
v1.19.4-v1.19.16
v1.20.1-v1.20.14
v1.21.1-v1.21.8
v1.22.1-v1.22.5
v1.23.1
```

## Usage
> k8s_ver=v1.23.1

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
kubeon add test \
  -m 192.168.60.26 \
  --master-name test60 \
  --default-passwd 123456 \
  --log-level debug
# add one worker
kubeon add test \
  -w 192.168.60.24 \
  --worker-name test40 \
  --default-passwd 123456 \
  --log-level debug
```

### Del node

```shell
# delon one node
kubeon del test \
    ip=192.168.60.24 \
    --log-level debug
# or use hostname
kubeon del test \
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
