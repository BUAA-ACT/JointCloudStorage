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
popd || exit
echo "Done"

popd || exit

exit 0