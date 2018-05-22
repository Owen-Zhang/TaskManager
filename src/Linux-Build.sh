#!/bin/bash

#编译golang文件

#if [ "$#" -ne 1 ]; then
#	echo "please input image tag"
#	exit
#fi

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/linux/home home/main.go

if [ $? -ne 0 ]; then
	echo "编译出错"
	exit
fi

echo "golang代码编译成功"

ver=$(cat version.txt)

echo ${ver}

docker build -t taskmanager:${ver} .

if [ $? -ne 0 ]; then
	echo "docker build 出现错误"
	exit
fi

echo "docker 打包成功"

docker save taskmanager:${ver} -o out/linux/taskmanager_${ver}.tar

if [ $? -ne 0 ]; then
	echo "打包镜像文件失败"
	exit
fi

echo "已生成了镜像tar文件在out/linux目录中"

#复制配制文件
rm -rf out/linux/conf
cp -R conf/  out/linux/conf

#删除可运行文件
rm -rf out/linux/home