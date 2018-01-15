package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/redis"
)

const lock_key = "LOCK_KEY"
const lock_value = "LOCK_VALUE"

func main() {
	client, _ := redis.Dial("tcp", "localhost:6379")
	defer client.Close()

	repl := client.Cmd("PING")
	content, _ := repl.Str()
	fmt.Println(content)
}

//我们首先使用GET方法，获取键的值，并把这个值转换为字符串，然后用if方法去检查有没有值，
//如果没有值的话就返回一个空的字符串，确认没有值就调用set方法进行设置，就是给它加锁
func acquire(client *redis.Client) bool {
	current_value, _ := client.Cmd("GET", lock_key).Str()
	if current_value == "" {
		client.Cmd("MULTI")
		client.Cmd("SET", lock_key, lock_value)
		rep, _ := client.Cmd("EXEC").List()
		if rep != nil {
			return true
		}
	}
	return false
}

// 使用NX方式加锁（获取锁）
//当我们使用带NX选项的SET命令时，只有在键key不存在的情况下才会对它进行设置，如果键已经有值，就会放弃对它进行设置代码，并返回nil表示设置失败
func acquireNX(client *redis.Client) bool {
	rep, _ := client.Cmd("SET", lock_key, lock_value, "NX").Str()
	return rep != ""
}

// 释放🔐就是 将lock_key del
func release(client *redis.Client) {
	client.Cmd("DEL", lock_key)
}

// 使用集合 进行在线统计,当一个用户上线的时候我们就把用户名添加在在线用户集合里面
const online_user_set = "ONLINE_USER_SET"

func set_online(client *redis.Client, user string) {
	client.Cmd("SADD", online_user_set, user) // 把用户名添加到集合
}

func count_online(client *redis.Client) int64 {
	rep, _ := client.Cmd("SCARD", online_user_set).Int64() // 获取集合元素数量
	return rep
}

// 为每一个用户创建一个ID，当一个用户上线后，就用他的ID作为索引，假设现在有一个用户peter，我们给他映射一个ID 10086，
// 然后根据这个ID把这个位图里面索引为10086的值设置为1，值为1的用户就是在线，值为0的就是不在线。
// 这里同样需要用到3个命令：
// SETBIT bitmap index value ：将位图指定索引上的二进制位设置为给定的值
// GETBIT bitmap index  ：获取位图指定索引上的二进制位
// BITCOUNT bitmap ：统计位图中值为 1 的二进制位的数量

const online_user_bitmap = "ONLINE_USER_BITMAP"

func set_online(client *redis.Client, user_id int64) {
	client.Cmd("SETBIT", online_user_bitmap, user_id, 1)
}

func count_online_bitmap(client *redis.Client) int64 {
	rep, _ := client.Cmd("BITCOUNT", online_user_bitmap).Int64()
	return rep
}

func is_online_or_not(client *redis.Client, user_id int64) bool {
	rep, _ := client.Cmd("GETBIT", online_user_bitmap, user_id).Int()
	return rep == 1
}

// 跟刚才的集合相比，虽然位图的体积仍然会随着用户数量的增多而变大，但因为记录每个用户所需的内存数量从原来的平均10字节变成了1位，
// 所以将节约大量的内存，把几十G的占用降为了几百M

// 我们要继续进行优化就得到了方法三——使用Hyperloglog。当一个用户上线时，我们就使用Hyperloglog对他进行计数。
// 假设现在有一个用户jack，我们通过Hyperloglog算法对他进行计数，然后把这个计数反映到Hyperloglog里面，
// 如果这个元素之前没有被Hyperloglog计数过的话，你新添加在Hyperloglog里面就会对自己的计数进行加1。如果jack已经存在，它的计数值就不会加1

const online_user_hll = "ONLINE_USER_HLL"

func set_online_hll(client *redis.Client, user string) {
	client.Cmd("PFADD", online_user_hll, user)
}

func count_online_hll(client *redis.Client) int64 {
	rep, _ := client.Cmd("PFCOUNT", online_user_hll).Int64()
	return rep
}

// redis 使用有序集合 自动补全功能
// 实现我们的自动补全需要用到两个命令，第一个ZINCRBY zset increment member是对给定成员的分值执行自增操作；
// 第二个ZREVRANGE zset start end [WITHSCORES]是按照分值从大到小的顺序，从有序集合里面获取指定索引范围内的成员。
// 因为我们的权重是从大到小排列

const autocomplete = "autocomplete::"

func feed(client *redis.Client, content string, weight int) {
	for i, _ := range content {
		segment := content[:i+1] // 枚举字符串组成排列
		key := autocomplete + segment
		client.Cmd("ZINCRBY", key, weight, content) // 对各个权重表进行更新
	}
}

func hint(client *redis.Client, prefix string, count int) []string {
	key := autocomplete + prefix
	result, _ := client.Cmd("ZREVRANGE", key, 0, count-1).List() // 按权重从大到小
	return result
}
