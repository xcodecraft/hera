ROOT_DIR=`pwd`
export GOPATH=$ROOT_DIR
export GOBIN=$GOPATH/bin
echo $GOPATH
echo $GOBIN

go install gopkg.in/yaml.v2
go install main
