#!/bin/bash

set -e

mysqlpass=111
volumename=yxi-back_api-db
containername=yxi-mariadb

docker volume create --name $volumename &&\
docker run --name $containername -d -e MYSQL_ROOT_PASSWORD=$mysqlpass -v $volumename:/var/lib/mysql -d mariadb:10.3 &&\
docker stop $containername &&\
docker rm $containername