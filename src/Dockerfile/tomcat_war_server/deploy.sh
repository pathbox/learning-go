#!/bin/sh

docker build -t my_tomcat_war/tomcat:v1 .

docker run -v /tmp/my_war_logs:/usr/local/tomcat/logs -d -p 8080:8080 my_tomcat_war/tomcat:v1 # 第二个端口号为docker中Tomcat war包服务的端口号，需要和server.xml中的配置的端口号一致，默认是8080 挂载的日志目录，docker运行用户要有读写和进入的权限

sleep 5 # 睡眠5秒等待Tomcat war完全启动

curl -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"test":{"id":"xxxxxx", "key":"xxxxxx","itsmApiUrl":"http://example.com"}}' \
"http://server:port/init" -v