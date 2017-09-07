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
	generator = func() string {
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
