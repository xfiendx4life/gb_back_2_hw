#!bin/sh
docker rm -f elasticsearch
docker run -d --name elasticsearch -p 9200:9200 -e discovery.type=single-node \
elasticsearch:7.17.3