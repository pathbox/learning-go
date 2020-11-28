该功能是手动graceful 重启http server
新版本编译好后，利用这种方式重启服务，能够避免老的connecttion 不会因为重启而中断，而是会继续处理完成

测试结果：

对于修改配置文件，平滑重启可以读取新的配置文件，或类似的全局变量数据

而重新编译的二进制程序，是无法使用的。fork的新进程的代码还是使用原有的



endless:

nice lib

能够实现重新编译的二进制程序， 能够用于平滑更新新的代码


### grace的实际操作例子描述
// curl "http://127.0.0.1:5001/sleep?duration=60s" &

/*
1.
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 62844 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)
启动服务后，看到服务的PID是62844

2. 执行请求，这个请求会持续60秒，这个请求连接的
curl "http://127.0.0.1:5001/sleep?duration=60s" &

netstat -anlt| grep 5001
tcp4       0      0  127.0.0.1.5001         127.0.0.1.61724        ESTABLISHED
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 62844 pathbox    4u  IPv6 0x4ce79a57f7870183      0t0  TCP localhost:commplex-link->localhost:61748 (ESTABLISHED) 可以看到这个是curl和服务建立的连接
example 62844 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)

3.
kill -USR2 62844

再执行：
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 62844 pathbox    4u  IPv6 0x4ce79a57f7870183      0t0  TCP localhost:commplex-link->localhost:61748 (ESTABLISHED)
curl    65640 pathbox    3u  IPv4 0x4ce79a5804a1230b      0t0  TCP localhost:61748->localhost:commplex-link (ESTABLISHED)
example 65870 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)
example 65870 pathbox    4u  IPv6 0x4ce79a57ef2f3183      0t0  TCP localhost:commplex-link->localhost:61779 (ESTABLISHED)
curl    67548 pathbox    3u  IPv4 0x4ce79a57f93cb153      0t0  TCP localhost:61779->localhost:commplex-link (ESTABLISHED)
kill -USR2已经执行了，看到新的LISTEN已经创建，PID是65870，但是curl 的连接还在一直还在执行，并没有因为kill操作而断了。
这时候再执行一个新的curl请求，新的请求：localhost:61779 连接是和新的服务进程PID65870 建立了连接，请的请求请求到了重启后的服务进程上，不会请求到旧的服务进程上。旧的服务进程等原有的请求执行完了，就释放了

4. 最后留下最新的服务进程
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 65870 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)

总结： grace重启后，旧的服务进程和新的服务进程同时存在，旧的服务进程在其原有的连接请求执行完后释放，新的请求会请求到新的服务进程进行处理

- 监听信号
- 收到信号时fork子进程（使用相同的启动命令），将服务监听的socket文件描述符传递给子进程
- 子进程监听父进程的socket，这个时候父进程和子进程都可以接收请求
- 子进程启动成功之后，父进程停止接收新的连接，等待旧连接处理完成（或超时）
- 父进程退出，升级完成


grace	Hello world	Hello Harry	2096	3100	旧API不会断掉，会执行原来的逻辑，pid会变化
endless	Hello world	Hello Harry	22072	22365	旧API不会断掉，会执行原来的逻辑，pid会变化
overseer	Hello world	Hello Harry	28300	28300	旧API不会断掉，会执行原来的逻辑，主进程pid不会变化
overseer是与grace和endless有些不同，主要是两点：

overseer添加了Fetcher，当Fetcher返回有效的二进位流(io.Reader) 时，主进程会将它保存到临时位置并验证它，替换当前的二进制文件并启动。
Fetcher运行在一个goroutine中，预先会配置好检查的间隔时间。Fetcher支持File、GitHub、HTTP和S3的方式。详细可查看包package fetcher
overseer添加了一个主进程管理平滑重启。子进程处理连接，能够保持主进程pid不变
*/
