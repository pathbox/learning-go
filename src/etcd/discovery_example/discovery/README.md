https://github.com/daizuozhuo/etcd-service-discovery

整个代码的思路很简单, worker启动时向etcd注册自己的信息,并设置一个过期时间TTL,每隔一段时间更新这个TTL,如果该worker挂掉了,这个TTL就会expire. master则监听workers/这个etcd directory, 根据检测到的不同action来增加, 更新, 或删除worker.

WatcherOptions里recursive指的是要监听这个文件夹下面所有节点的变化, 而不是这个文件夹的变化. 当返回expire的时候, 该节点不一定挂掉, 有可能只是网络状况不好, 因此我们只将它暂时设置成不在集群里, 等当它返回update时在设置回来. 只有返回delete才明确表示将它删除

worker这边也跟master类似, 保存一个etcd KeysAPI, 通过它与etcd交互.然后用heartbeat来保持自己的状态

每个Worker 可以看成是一个 Server

Master 看成是 etcd的 proxy, 也就是注册中心

例子中,是Worker 不断的给 Master 发心跳,以维持联系. 在TTL的时间内,如果没有发心跳,进行set刷新TTL,key就会在etcd中删除.
也就表示这个Server 挂了

还有一种方式,刚好是相反的. 就是Master 调所有Worker, 如果调通,表示 worker是活的,调失败,表示 worker挂了,这样一定程度会增加etcd的负载.

个人会选择第一种服务发现的模式.

注册中心接受 Server的心跳,在一定程度上,可以将负载分配到各个Server上