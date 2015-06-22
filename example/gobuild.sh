#!/bin/bash
ROOT_DIR=`pwd`

export GOPATH=$ROOT_DIR
export GOBIN=$ROOT_DIR/bin

echo "Init env..."
echo "GOPATH is :$GOPATH" 
echo "GOBIN  is :$GOBIN" 

echo "It is deploying hera, please wait..."
go get -u github.com/xcodecraft/hera
echo "hera deploy success."

rm -rf ./pkg/ ./bin/

echo "It is complie your prj, please wait..."
go install main
echo "complie success."

exec ./bin/main
