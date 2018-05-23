#!/bin/sh

go install
rm -rf ./bin/kk-room
cp $GOPATH/bin/kk-room ./bin/kk-room
