#!/bin/bash
trap "rm server;kill 0" EXIT


go build -o server
./server -port=8001 &
./server -port=8002 &

#启动了API服务和Cache服务
./server -port=8003 -api=1 &

sleep 2
echo ">>> Server started ,now let's begin test <<<"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait
