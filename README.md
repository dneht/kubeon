## Install

> sh -c "$(wget https://dl.sre.pub/on/install.sh -q -O -)"

## Support

> kubeon view support

```text
v1.19.4-v1.19.11
v1.19.7-v1.20.1
v1.21.1
```

## Usage

### Vagrant test

> cd test/ubuntu20 && vagrant up

### Create cluster

```shell
# ssh node0
vagrant ssh node0
# use root with password 4567890123
su - root
# install kubeon
sh -c "$(wget https://dl.sre.pub/on/install.sh -q -O -)"
# create cluster
kubeon create -N test --version v1.21.1 \
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
Please use `source /etc/profile` for auto completion

### -N
cluster name

#### --cri
default is `ontainerd`, you can also use `docker`

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

### Add node

```shell
# add one master
kubeon add -N test \
    -m 172.20.0.24 \
    --master-name test40 \
    --default-passwd 4567890123 \
    --log-level debug
# add one worker
kubeon add -N test \
    -w 172.20.0.26 \
    --worker-name test110 \
    --default-passwd 4567890123 \
    --log-level debug
```

### Del node

```shell
# del one node
kubeon del -N test \
    ip=172.20.0.24 \
    --log-level debug
# or use hostname
kubeon del -N test \
    name=test40 \
    --log-level debug
```

### Upgrade cluster

```shell
kubeon upgrade -N test --version v1.21.1 \
    --log-level debug
```

### Destroy cluster

```shell
kubeon destroy -N test \
    --log-level debug
```

### Cluster info

```shell
# cluster info
kubeon view cluster-info -N test
# all node ipvs rule
kubeon exec @all "ipvsadm -ln" -N test -R
```