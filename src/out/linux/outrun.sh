#!/bin/bash

if [ "$#" -ne 2 ]; then
       echo "please input image tag and port, first image tag, second port"
       exit
fi


#需要两个参数:
#1:镜像版本号 taskmanager:$1
tag=$1

#2:本机开放的端口号$2:8000
port=$2

if [ -f "/opt/taskmanager/conf/app.conf" ]; then
        rm -rf /opt/taskmanager/conf/app.conf
else
        mkdir -p /opt/taskmanager/conf
fi

if [ -d "/opt/taskmanager/logs/" ]; then
        rm -rf /opt/taskmanager/logs
fi

\cp conf/app.conf /opt/taskmanager/conf/

#删除运行中的容器
containerid=`docker ps -a | grep taskmanager | awk '{print $1}'`
docker rm -f ${containerid}

#删除老版本的镜像
imageid=`docker images | grep taskmanager | awk '{print $3}'`
docker rmi ${imageid}

#加载镜像
docker load -i taskmanager_${tag}.tar

#删除以前版本的镜像(这个要通过 docker ps -a | grep taskmanager 得出结果去获得相应的containerId)

docker run -d -p ${port}:8000 --name=taskmanager -v /opt/taskmanager/logs:/logs -v /opt/taskmanager/conf:/conf taskmanager:${tag}