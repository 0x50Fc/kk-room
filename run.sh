#!/bin/sh

start() {
	nohup ./kk-room -admin :8081 -app :8082 -broadcast :8080 > kk-room.log 2>&1 &
	echo "$!" > .pid
}

stop() {
	kill `cat .pid`
}

restart() {
	stop
	start
}

if [[ $1 == "start" ]] ; then
        start
elif [[ $1 == "stop" ]] ; then
        stop
else
        restart
fi
