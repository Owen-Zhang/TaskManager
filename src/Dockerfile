#FROM frolvlad/alpine-glibc:alpine-3.7 要用到shell相关,此镜像没法满足
FROM centos:7
MAINTAINER Owen.zhang

add out/linux/home  /
add conf/  /conf
add views/ /views
add static/  /static
add logs/   /logs

RUN mkdir -p /Data/Temp \
 && mkdir -p /Data/Run

#RUN apk add --no-cache tzdata \
#   &&  /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
#   && echo 'Asia/Shanghai' >/etc/timezone
# 时间通过外面的-v传入 -v /etc/localtime:/etc/localtime

EXPOSE 8000
CMD ["/home"] 