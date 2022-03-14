package ratelimiter

import (
	"context"
	"fmt"
	"time"
)

func RateAllow() {
	limiter := rate.NewLimiter(10, 100)

	for i := 0; i < 20; i++ {
		if limiter.AllowN(time.Now(), 25) {
			// TODO 业务
		} else {
			// limit
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func rateWait() {
	limiter := rate.NewLimiter(1, 10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for i := 0; ; i++ {
		err := limiter.WaitN(ctx, 2) // 获取不到会阻塞或者超时
		if err != nil {
			return
		}
		// TODO 业务
	}
	fmt.Println("done")
}

func rateReserve() {
	limiter := rate.NewLimiter(3, 5)

	// 动态修改桶大小
	limiter.SetBurst(100)
	// 动态修改生成令牌速率
	// limiter.SetLimit(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for i := 0; ; i++ {
		fmt.Printf("%3d %s\n", i, time.Now().Format("2006-01-02 15:04:05.000"))
		reserveN := limiter.ReserveN(time.Now(), 4)
		if !reserveN.OK() {
			// 返回异常
			// fmt.Println("Not allowed to act! Did you remember to set lim.burst to be  0?")
		}
		delay := reserveN.Delay()
		fmt.Println("sleep delay ", delay)
		time.Sleep(delay)

		select {
		case <-ctx.Done():
			fmt.Println("timeout, quit")
			return
		default:
		}
		// TODO 业务
	}

	fmt.Println("main")
}
