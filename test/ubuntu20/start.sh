num=$1
ver=v1.19.11
pwd=4567890123
echo "start node ${num} with password: ${pwd}"
sudo echo root:${pwd} | chpasswd
sudo sed -i "/PermitRootLogin prohibit-password/d" /etc/ssh/sshd_config
sudo sed -i "s/PasswordAuthentication no/PasswordAuthentication yes/g" /etc/ssh/sshd_config
sudo echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
sudo systemctl restart sshd
sudo echo "apt_preserve_sources_list: true" >> /etc/cloud/cloud.cfg
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal main restricted" > /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted" >> /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal universe" >> /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-updates universe" >> /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal multiverse" >> /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-updates multiverse" >> /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted" >> /etc/apt/sources.list
sudo echo "deb http://mirrors.aliyun.com/ubuntu/ focal-security universe" >> /etc/apt/sources.list
sudo apt-get update
sudo apt-get upgrade
sudo apt-get install -y chrony
sudo systemctl start chrony
sudo sh -c "$(wget https://dl.sre.pub/on/install.sh -q -O -)"
if [ $num = 6 ]; then
  kubeon create -N test --version ${ver} \
      -m 172.20.0.21 \
      -m 172.20.0.22 \
      -m 172.20.0.23 \
      --master-name test10 \
      --master-name test20 \
      --master-name test30 \
      -w 172.20.0.24 \
      -w 172.20.0.25 \
      --worker-name test40 \
      --worker-name test50 \
      --default-passwd ${pwd} \
      --interface enp0s8 \
      --log-level debug
  kubeon view cluster-info -N test
  sleep 2s
  kubeon add -N test \
      -m 172.20.0.26 \
      --master-name test60 \
      --default-passwd ${pwd} \
      --log-level debug
  kubeon view cluster-info -N test
  sleep 2s
  kubeon del -N test \
      ip=172.20.0.25 \
      --log-level debug
  kubeon view cluster-info -N test
  sleep 2s
fi