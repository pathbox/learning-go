SET key value [EX seconds] [PX milliseconds] [NX|XX]

加锁操作

SET user:id:1 1 PX 100 NX  // 锁的有效时间这里为 100 milliseconds，根据实际场景设置

如果在有效时间内，想要锁的 业务逻辑代码没有执行完，锁也会释放，这是锁的自动释放时间，其他线程也可以获得锁，这时候就会出现问题



关于redis分布式锁的讨论

https://redis.io/topics/distlock

http://antirez.com/news/101

http://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html

https://huoding.com/2015/09/14/463

如果一个请求更新缓存的时间比较长，甚至比锁的有效期还要长，导致在缓存更新过程中，锁就失效了，此时另一个请求会获取锁，当前一个请求在缓存更新完毕之后，如果不加以判断直接删除锁，就会出现误删除其它请求创建的锁的情况，所以在创建所得时候，需要引入一个随机唯一值

```go
status, err := redis.String(lock.conn.Do("SET", lock.key(), lock.token, "EX", int64(lock.timeout/time.Second)))

if status == "OK" {
  cache.update()
  if (redis.Get(lock.key()) == lock.token) {
    redis.Del(logk.key(j))
  }
}

```