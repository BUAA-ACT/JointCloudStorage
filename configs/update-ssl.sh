#!/bin/bash

# export Ali_Key="LTAI4G3**********"
# export Ali_Secret="5bmnIvUqvu*********************"
acme.sh --issue --dns dns_ali -d jointcloudstorage.cn -d *.jointcloudstorage.cn --dnssleep 30
acme.sh --install-cert -d jointcloudstorage.cn --key-file /home/jcspan/ssl/jointcloudstorage.cn/key.pem --fullchain-file /home/jcspan/ssl/jointcloudstorage.cn/cert.pem --reloadcmd "sudo service nginx force-reload"

host_list={"aliyun-hohhot.jointcloudstorage.cn", "aliyun-qingdao.jointcloudstorage.cn",
"aliyun-hangzhou.jointcloudstorage.cn", "txyun-chengdu.jointcloudstorage.cn",
"bdyun-guangzhou.jointcloudstorage.cn"}

for host in ${host_list[@]}; do
    echo "Install SSL for $host"
    scp -r /home/jcspan/ssl/jointcloudstorage.cn/ jcspan@$host:/home/jcspan/ssl/
    scp ./nginx-host.conf jcspan@$host:/etc/nginx/conf.d/jcspan.conf
    ssh -t $host "sudo service nginx restart"
    echo "Install SSL for $host success"
done

