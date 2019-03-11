// https://mp.weixin.qq.com/s?__biz=MzI1MTIzMzI2MA==&mid=2650561116&idx=1&sn=a6cb8c7bf52bc94b1a5d9feae21effa2&chksm=f1feecdfc68965c9bf20c1ef9373118dcf79fdc93a29dd172796b0dbef3e9ef6830b44ef072f&mpshare=1&scene=24&srcid=0616KUfpBJBFKKcCh0Z69GtW&key=b10a7c153a57fbb9c9340e20573980776f826657ab7cdb9b076ed3ca96bcd772a26f9f91086f829947b1285734530a1356d6fd2476485bba6da9f8f4f4637fc5ac69a884093e6b5ceb391093e924d291&ascene=0&uin=OTcxOTY1NTU%3D&devicetype=iMac+MacBookPro12%2C1+OSX+OSX+10.12+build(16A323)&version=12020710&nettype=WIFI&fontScale=100&pass_ticket=mmYmEVv3gqNbe2uX0CV7S0tVNyYKDKJ9qiaR9Jf5%2Fno%3D

// 时间负责度 O(N) 空间复杂度 O(1)
package main

import (
	"fmt"
)

func main() {
	n := 10
	result := getClimbingWays(n)
	fmt.Println("The result is: ", result)
}

func getClimbingWays(n int) int {
	if n < 1 {
		return 0
	} else if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}

	a := 1
	b := 2
	tmp := 0

	for i := 3; i <= n; i++ {
		tmp = a + b
		a = b
		b = tmp
	}
	return tmp
}
