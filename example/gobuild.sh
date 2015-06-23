#!/bin/bash
ROOT_DIR=`pwd`

export GOPATH=$ROOT_DIR
export GOBIN=$ROOT_DIR/bin

echo "Init env..."
echo "GOPATH is :$GOPATH"
echo "GOBIN  is :$GOBIN"

echo "It is deploying hera, please wait..."
go get github.com/xcodecraft/hera
echo "hera deploy success."

rm -rf ./pkg/ ./bin/

echo "It is complie your prj, please wait..."
go install main
echo "complie success."

sudo rm -rf /usr/local/nginx/conf/include/example_duanbingying_nginx.conf
sudo ln -s $ROOT_DIR/conf/nginx.conf /usr/local/nginx/conf/include/example_duanbingying_nginx.conf

sudo nginx -s reload

/usr/local/bin/python  /usr/local/bin/zdaemon  -C $ROOT_DIR/conf/zd.xml start
