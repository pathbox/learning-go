FROM tomcat
MAINTAINER "pathbox"
ADD my_server.war /usr/local/tomcat/webapps/
# 修改了好的server.xml文件
COPY server.xml /usr/local/tomcat/conf/
CMD ["catalina.sh", "run"]

