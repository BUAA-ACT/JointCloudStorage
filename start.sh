#!/bin/bash

base_dir=$(pwd)
build_dir="$base_dir/build/"
configs_dir="$base_dir/configs/private/"

if [ ! -d "$configs_dir" ]; then
    exit 1
fi

pushd "$configs_dir" || exit
docker-compose down
docker-compose up -d 
echo "MongoDB 启动成功"
popd || exit

screen_do(){
    screen_name=$1
    cmd=$2
    if screen -ls "$screen_name" > /dev/null ; then 
        screen -S "$screen_name" -X quit
    fi
    echo "创建 screen: $screen_name"
    screen -dmS "$screen_name"
    screen -x -S "$screen_name" -p 0 -X stuff "$cmd"
    screen -x -S "$screen_name" -p 0 -X stuff $'\n' 
}

echo "等待 mongoDB 启动"
sleep 5

pushd "$build_dir" || exit
screen_do "aliyun-hohhot-h" "./httpserver -c ../configs/private/0/httpserver.properties "
screen_do "aliyun-hohhot-t" "./transporter -c ../configs/private/0/transporter_config.json "
screen_do "aliyun-hohhot-s" "./scheduler -addr=:8082 -cid=aliyun-hohhot  -env=dev -heartbeat=10s -reschedule=60s "


popd || exit


