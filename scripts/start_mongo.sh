#!/bin/bash

set -e 

docker rm mongo -f 2> /dev/null || true

docker run --rm -d -p 27017:27017 -h localhost --name mongo mongo:5.0.9 --replSet=rs0
sleep 4
docker exec mongo mongo --eval 'rs.initiate({_id:"rs0", members:[ {_id:0, host:"127.0.0.1:27017"}]});'
