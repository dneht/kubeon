num=$1
ver=v1.28.1
pwd=123456
echo "start node ${num} with password: ${pwd}"
sudo echo root:${pwd} | chpasswd
sudo sed -i "/PermitRootLogin prohibit-password/d" /etc/ssh/sshd_config
sudo sed -i "s/PasswordAuthentication no/PasswordAuthentication yes/g" /etc/ssh/sshd_config
sudo echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
sudo systemctl restart sshd
export DEBIAN_FRONTEND=noninteractive
sudo echo "apt_preserve_sources_list: true" >> /etc/cloud/cloud.cfg
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal main restricted" > /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted" >> /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal universe" >> /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-updates universe" >> /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal multiverse" >> /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-updates multiverse" >> /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted" >> /etc/apt/sources.list
#sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-security universe" >> /etc/apt/sources.list
sudo apt-get update
sudo apt-get install -y chrony
sudo systemctl start chrony
sudo sh -c "$(wget https://back.pub/kubeon/install.sh -q -O -)"
if [ $num = 5 ]; then
  kubeon create test ${ver} \
      --cilium-enable-dsr \
      -m 192.168.60.21 \
      -m 192.168.60.22 \
      -m 192.168.60.23 \
      --master-name test10 \
      --master-name test20 \
      --master-name test30 \
      -w 192.168.60.24 \
      --worker-name test40 \
      --default-passwd ${pwd} \
      --interface enp0s8 \
      --ic contour \
      --with-kata \
      --with-kruise \
      --offline \
      --v 6
  sleep 2s
  kubeon display test
  kubeon add test \
      -m 192.168.60.25 \
      --master-name test50 \
      --default-passwd ${pwd} \
      --v 6
  sleep 4s
  kubeon display test
  kubeon del test \
      ip=192.168.60.25 \
      --v 6
  sleep 4s
  kubeon display test

  cat>"${HOME}/test.yaml"<<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
  labels:
    run: app
    type: test
spec:
  replicas: 3
  selector:
    matchLabels:
      run: app
      type: test
  template:
    metadata:
      labels:
        run: app
        type: test
    spec:
      containers:
        - name: test
          image: registry.cn-hangzhou.aliyuncs.com/dneht/debian-test:latest
          args:
            - /bin/sh
            - -c
            - sleep 10; touch /tmp/healthy; sleep 30000
          readinessProbe:
            exec:
              command:
                - cat
                - /tmp/healthy
            initialDelaySeconds: 10
            periodSeconds: 5
EOF
  kubectl apply -f ${HOME}/test.yaml
fi
