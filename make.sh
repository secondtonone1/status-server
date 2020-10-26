#!/bin/bash -x

if [ ! -f "$0" ]; then
    echo 'make must be run within its container folder' 1>&2
    exit 1
fi


if [ "$1" == "clean" ]; then
        rm bin/* pkg/* utility/pkg/* -rf
        exit
fi

CURDIR=`pwd`
export GOPATH="$CURDIR"
export GOBIN=${GOPATH}/bin
author=${USER}
date=`date "+%Y-%m-%d_%H:%I:%S"`

goversion=`go version`

ldflags="-X status-server/logic._AUTHOR_=$author -X status-server/logic._COMPILETIME_=\"$date\""
echo ldflags ${ldflags}

echo "formating code..."
gofmt -w src/

# 以下命令可以使用go get golang.org/x/tools/cmd/goimports获取
#${GOROOT}/bin/goimports -w=true src/

cd ./src/
#${GOROOT}/bin/go install -v -ldflags "$ldflags"  main/main.go
#go install -v -ldflags "$ldflags"  main/main.go
CGO_ENABLED=0 GOOS=linux go build -v -ldflags "$ldflags" -a -installsuffix cgo -o status-server  main/main.go
mv status-server ../bin/

if [ $? == 0 ]; then
	#mv ../bin/main ../bin/status-server
	echo "build success"
else
	echo "build error"
fi
