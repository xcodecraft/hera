ROOT_DIR=`pwd`
export GOPATH=$ROOT_DIR
export GOBIN=$GOPATH/bin
echo $GOPATH
echo $GOBIN

go install main
