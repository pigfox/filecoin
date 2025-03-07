#!/bin/sh
clear
set -x
set -e

docker-compose down -v
docker image prune -a -f
docker-compose build --no-cache
docker-compose up