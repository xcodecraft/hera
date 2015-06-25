#!/bin/bash
ROOT_DIR=`pwd`

PRJ='hera_example'

export GOPATH=$ROOT_DIR
export GOBIN=$ROOT_DIR/bin

if [ $1 = "clean" ]; then
   echo "clean..."
   rm -rf ./pkg ./bin ./src/github.com/
   echo "clean finish"
   exit
fi

echo "**********************************Strat****************************"

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

case $1 in
    start)
        echo "$PRJ start:"
        sudo rm -rf /usr/local/nginx/conf/include/$PRJ.conf
        sudo ln -s $ROOT_DIR/conf/nginx.conf /usr/local/nginx/conf/include/$PRJ.conf
        sudo nginx -s reload
        /usr/local/bin/python  /usr/local/bin/zdaemon  -C $ROOT_DIR/conf/zd.xml start
        ;;
    stop)
        echo "$PRJ stop:"
        /usr/local/bin/python  /usr/local/bin/zdaemon  -C $ROOT_DIR/conf/zd.xml stop
        ;;
    restart)
        echo "$PRJ restart:"
        sudo rm -rf /usr/local/nginx/conf/include/$PRJ.conf
        sudo ln -s $ROOT_DIR/conf/nginx.conf /usr/local/nginx/conf/include/$PRJ.conf
        sudo nginx -s reload
        /usr/local/bin/python  /usr/local/bin/zdaemon  -C $ROOT_DIR/conf/zd.xml restart
        ;;
esac
echo "**********************************End****************************"
