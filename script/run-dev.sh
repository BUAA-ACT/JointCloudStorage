#!/bin/bash

base_dir=$(pwd)
build_dir="$base_dir/build/"

if [ ! -d "$base_dir/build/" ]; then
    mkdir "$base_dir/build/"
fi
pushd "$base_dir/service" || exit
echo "编译 transporter"
pushd "transporter" || exit
make build
cp build/bin/transporter $build_dir
popd || exit
echo "Done"

echo "编译 scheduler"
pushd "scheduler" || exit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s -extldflags "-static"' -o build/bin/scheduler
cp build/bin/scheduler $build_dir
popd || exit
echo "Done"

echo "编译 httpserver"
pushd "httpserver" || exit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s -extldflags "-static"' -o build/bin/httpserver
cp build/bin/httpserver $build_dir
cp httpserver.properties $build_dir
popd || exit
echo "Done"

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

pushd "$build_dir" || exit
screen_do "aliyun-hohhot-h" "./httpserver"
screen_do "aliyun-hohhot-s" "./scheduler -addr=:8082 -cid=aliyun-hohhot  -env=aliyun-hohhot -heartbeat=10s -reschedule=60s -mongo=mongodb://localhost:27017"

screen_do "aliyun-qingdao-s" "./scheduler -addr=:8282 -cid=aliyun-qingdao -env=aliyun-qingdao -heartbeat=10s -reschedule=60s -mongo=mongodb://localhost:27017"

screen_do "aliyun-hangzhou-s" "./scheduler -addr=:8182 -cid=aliyun-hangzhou -env=aliyun-hangzhou -heartbeat=10s -reschedule=60s -mongo=mongodb://localhost:27017"




