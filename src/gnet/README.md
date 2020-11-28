https://gnet.host/blog/presenting-gnet/

https://github.com/panjf2000/gnet

>gnet is not designed to displace the standard Go net package, but to create a networking server framework for Go that performs on par with Redis and Haproxy for networking packets handling.

不是为了取代Go的 net 库，而是为某些特定场景提供性能更好的功能.
对于一般web api server, Go 的net库已经封装的很好，实现的很好。

`gnet`在一些功能和协议上没有封装，偏向更底层tcp传输。你操作的都是byte数据，而不是文本数据