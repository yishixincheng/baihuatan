#!/bin/bash

firewall-cmd --add-port=2479/tcp --permanent --zone=public
firewall-cmd --add-port=2380/tcp --permanent --zone=public
firewall-cmd --add-port=2480/tcp --permanent --zone=public
firewall-cmd --add-port=2381/tcp --permanent --zone=public

current_file_path=$(cd "$(dirname "$0")"; pwd)
cd ${current_file_path}

ETCD_INITIAL_CLUSTER="etcd0=http://192.168.43.251:2380,etcd1=http://192.168.43.251:2381"

ETCD_INITIAL_CLUSTER_STATE=new

#export currentHostIp=`ip -4 address show eth0 | grep 'inet' |  grep -v grep | awk '{print $2}' | cut -d '/' -f1`

firewall-cmd --reload
firewall-cmd --list-all

# node0

docker stop etcd0
docker rm   etcd0

rm -rf /tmp/etcd-data0.tmp && mkdir -p /tmp/etcd-data0.tmp && \
docker run -d \
--restart=always \
--hostname=etcd0 \
-p 2380:2380 \
-p 2479:2379 \
-v /etc/localtime:/etc/localtime \
--name etcd0 etcd \
/usr/local/bin/etcd -name etcd0 \
--data-dir /etcd-data \
--advertise-client-urls http://192.168.43.251:2479 \
--listen-client-urls http://0.0.0.0:2379 \
--initial-advertise-peer-urls http://192.168.43.251:2380 \
--listen-peer-urls http://0.0.0.0:2380 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster "etcd0=http://192.168.43.251:2380,etcd1=http://192.168.43.251:2381" \
--initial-cluster-state new \
--log-level info \
--logger zap \
--log-outputs stderr

#node1

docker stop etcd1
docker rm   etcd1

rm -rf /tmp/etcd-data1.tmp && mkdir -p /tmp/etcd-data1.tmp && \
docker run -d \
--restart=always \
--hostname=etcd1 \
-p 2381:2380 \
-p 2480:2379 \
-v /etc/localtime:/etc/localtime \
--name etcd1 etcd \
/usr/local/bin/etcd -name etcd1 \
--data-dir /etcd-data \
--advertise-client-urls http://192.168.43.251:2480  \
--listen-client-urls http://0.0.0.0:2379 \
--initial-advertise-peer-urls http://192.168.43.251:2381 \
--listen-peer-urls http://0.0.0.0:2380  \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster "etcd0=http://192.168.43.251:2380,etcd1=http://192.168.43.251:2381" \
--initial-cluster-state new \
--log-level info \
--logger zap \
--log-outputs stderr 