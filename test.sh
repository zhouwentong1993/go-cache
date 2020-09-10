#!/bin/zsh

trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 &
./server -port=8002 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9099/api?key=zhangsan1" &
curl "http://localhost:9099/api?key=zhangsan1" &
curl "http://localhost:9099/api?key=zhangsan1" &

sleep 100
