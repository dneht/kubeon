num=$1
ver=v1.22.1
pwd=123456
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
  kubeon create test ${ver} \
      -m 192.168.60.21 \
      -m 192.168.60.22 \
      -m 192.168.60.23 \
      --master-name test10 \
      --master-name test20 \
      --master-name test30 \
      -w 192.168.60.24 \
      -w 192.168.60.25 \
      --worker-name test40 \
      --worker-name test50 \
      --default-passwd ${pwd} \
      --ic contour \
      --interface enp0s8 \
      --log-level debug
  kubeon view info test
  sleep 2s
  kubeon add test \
      -m 192.168.60.26 \
      --master-name test60 \
      --default-passwd ${pwd} \
      --log-level debug
  kubeon view info test
  sleep 2s
  kubeon del test \
      ip=192.168.60.25 \
      --log-level debug
  kubeon view info test
  sleep 2s
fi