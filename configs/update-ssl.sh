#!/bin/bash

# export Ali_Key="LTAI4G3**********"
# export Ali_Secret="5bmnIvUqvu*********************"

sudopsw="jcspanisgreat"

source ~/.bashrc
shopt -s expand_aliases
mkdir -p  ~/ssl/jointcloudstorage.cn
# ~/.acme.sh/acme.sh --issue --dns dns_ali -d jointcloudstorage.cn -d *.jointcloudstorage.cn --dnssleep 30
~/.acme.sh/acme.sh --install-cert -d jointcloudstorage.cn --key-file ~/ssl/jointcloudstorage.cn/key.pem --fullchain-file ~/ssl/jointcloudstorage.cn/cert.pem 

host_list=("aliyun-hohhot.jointcloudstorage.cn" "aliyun-qingdao.jointcloudstorage.cn" 
"aliyun-hangzhou.jointcloudstorage.cn" "txyun-chengdu.jointcloudstorage.cn" 
"bdyun-guangzhou.jointcloudstorage.cn")

for host in ${host_list[@]}; do
    echo "Install SSL for $host"
    scp  ~/ssl/jointcloudstorage.cn/cert.pem jcspan@$host:/home/jcspan/ssl/cert.pem
    scp  ~/ssl/jointcloudstorage.cn/key.pem jcspan@$host:/home/jcspan/ssl/key.pem
    scp ./nginx-host.conf jcspan@$host:/home/jcspan/jcspan/nginx-host.conf
    ssh -tt jcspan@$host "sed -i 's/<host-name>/$host/g' /home/jcspan/jcspan/nginx-host.conf"
    echo $sudopsw | ssh -tt jcspan@$host "sudo cp /etc/nginx/conf.d/jcspan.conf /etc/nginx/conf.d/jcspan.conf.bak"
    echo $sudopsw | ssh -tt jcspan@$host "sudo cp /home/jcspan/jcspan/nginx-host.conf /etc/nginx/conf.d/jcspan.conf"
    echo $sudopsw | ssh -tt jcspan@$host "sudo service nginx restart"
    echo "Install SSL for $host success"
done

