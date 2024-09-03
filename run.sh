#!/bin/bash
#这里可替换为你自己的执行程序，其他代码无需更改
CMS_NAME=netflow

chmod 777 netflow
#使用说明，用来提示输入参数
usage() {
 echo "Usage: sh rundev.sh [start|stop|restart|status]"
 exit 1
}

#检查程序是否在运行
is_exist(){
 # shellcheck disable=SC2006
 # shellcheck disable=SC2009
 pid=`ps -ef|grep netflow|grep -v grep|awk '{print $2}' `
 #如果不存在返回1，存在返回0
 if [ -z "${pid}" ]; then
 return 1
 else
 return 0
 fi
}

#启动方法
start(){
 is_exist
 if [ $? -eq "0" ]; then
 echo "${CMS_NAME} is already running. pid=${pid} ."
 else

 nohup ./$CMS_NAME> /opt/www/server/logFileName1.file 2>&1 &
 echo "${CMS_NAME} start success"
 fi
}

#停止方法
stop(){
 is_exist
 if [ $? -eq "0" ]; then
 kill -9 $pid
 else
 echo "${CMS_NAME} is not running"
 fi
}

#输出运行状态
status(){
 is_exist
 if [ $? -eq "0" ]; then
 echo "${CMS_NAME} is running. Pid is ${pid}"
 else
 echo "${CMS_NAME} is NOT running."
 fi
}

#重启
restart(){
 stop
 start
}

#根据输入参数，选择执行对应方法，不输入则执行使用说明
case "$1" in
 "start")
 start
 ;;
 "stop")
 stop
 ;;
 "status")
 status
 ;;
 "restart")
 restart
 ;;
 *)
 usage
 ;;
esac