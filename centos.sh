#!/bin/sh

HOME=`pwd`

docker run --rm -v $GOPATH:/go:rw  -w /go/src/github.com/hailongz/kk-room --entrypoint ./build.sh registry.cn-beijing.aliyuncs.com/kk/kk-game:latest

rm -rf $HOME/../bin/kk-room
cp $HOME/bin/kk-room $HOME/../bin/kk-room
