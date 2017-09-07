// 加了订单数进去，做为权重来为用户抽奖。此题和上面的问题如此的相似，可把上面的问题， 理解成所有的用户权重都相同的抽奖，而此题是权重不同的抽奖。解决此问题，依旧是把map转为数组来思考， 把各用户的权重，从前到后依次拼接到数轴上，数轴的起点到终点即时中奖区间，而随机数落到的那个用户的区间，那个用户即为中奖用户

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	var users map[string]int64 = map[string]int64{
		"a": 10,
		"b": 5,
		"c": 15,
		"d": 20,
		"e": 10,
		"f": 30,
	}

	rand.Seed(time.Now().Unix())
	award_stat := make(map[string]int64)
	generator := GetAwardGenerator(users)
	for i := 0; i < 100000; i += 1 {
		name := generator()
		if count, ok := award_stat[name]; ok {
			award_stat[name] = count + 1
		} else {
			award_stat[name] = 1
		}
	}

	for name, count := range award_stat {
		fmt.Printf("user: %s, award count: %d\n", name, count)
	}
}

func GetAwardGenerator(users map[string]int64) (generator func() string) {
	var sum_num int64
	name_arr := make([]string, len(users))
	for u_name, num := range users {
		sum_num += num
		name_arr = append(name_arr, u_name)
	}
	generator = func() string { // 一个闭包
		award_num := rand.Int63n(sum_num)

		var offset_num int64
		for _, u_name := range name_arr {
			offset_num += users[u_name]
			if award_num < offset_num {
				return u_name
			}
		}

		return name_arr[0]
	}
	return
}

/*
GetAwardGenerator 方法会执行一次，这样，得到的name_arr 是固定的。也就是 抽象的 数轴区域是固定的。 这个数轴区域是权重之和，在这里就是 90。
在抽奖的方法中，闭包方法generator是循环执行。
随机得到一个抽奖数 award_num，把这个抽奖数抽象成为 抽奖数轴区域中的某个点 J(award_num, 0)

抽奖数轴区域的一种情况：（由于map是无序的，所以，实际在users数组中的顺序不是固定的）

a (0, 0)-(10, 0)
b (10, 0)-(15, 0)
c (15, 0)-(30, 0)
d (30, 0)-(50, 0)
e (50, 0)-(60, 0)
f (60, 0)-(90, 0)

f区域所占的权重是最大的。 上面方法使用了offset_num 偏移方法是以时间换空间。每一次抽奖都要range name_arr 这个数组

还有一种以空间换时间的方法，就是先构造一个 len(sum_num)的数组。数组中的值为不同区域的name。
举个例子就是 name_arr[60:90] 区域的数组的元素都是 f。这样，只要award_num随机抽取出来了，就知道
是谁中奖了。 name_arr[award_num] 再进行count计数。 就不会每次抽奖都要range 循环 name_arr，但是这里用了更多的内存存储空间，而不是一个变量存储空间。
*/
