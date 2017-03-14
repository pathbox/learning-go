package main

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type otherContext struct {
	context.Context
}

func main() {

	// 50 millisecond 之后自动关闭
	c, _ := context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
	if got, prefix := fmt.Sprint(c), "context.Background.WithDeadline("; !strings.HasPrefix(got, prefix) {
		fmt.Errorf("c.String() = %q want prefix %q", got, prefix)
	}
	testDeadline(c, "WithDeadline", time.Second)

	// 继承context的类型使用WithDeadline
	c, _ = context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
	o := otherContext{c}
	testDeadline(o, "WithDeadline+otherContext", time.Second)

	// Parent Context 超时时间比 Child Context 超时时间短
	c, _ = context.WithDeadline(context.Background(), time.Now().Add(50*time.Millisecond))
	o = otherContext{c}
	c, _ = context.WithDeadline(o, time.Now().Add(4*time.Second))
	// 结果为：Parent Context 超时
	testDeadline(c, "WithDeadline+otherContext+WithDeadline", 2*time.Second)

	// 超时时间是负数，直接超时
	c, _ = context.WithDeadline(context.Background(), time.Now().Add(-time.Millisecond))
	testDeadline(c, "WithDeadline+inthepast", time.Second)

	// 超时时间是现在时间，直接超时
	c, _ = context.WithDeadline(context.Background(), time.Now())
	testDeadline(c, "WithDeadline+now", time.Second)

	// WithTimeout

	c, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	if got, prefix := fmt.Sprint(c), "context.Background.WithDeadline("; !strings.HasPrefix(got, prefix) {
		fmt.Errorf("c.String() = %q want prefix %q", got, prefix)
	}
	testDeadline(c, "WithTimeout", time.Second)

	c, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	o = otherContext{c}
	testDeadline(o, "WithTimeout+otherContext", time.Second)

	c, _ = context.WithTimeout(context.Background(), 50*time.Millisecond)
	o = otherContext{c}
	c, _ = context.WithTimeout(o, 3*time.Second)
	testDeadline(c, "WithTimeout+otherContext+WithTimeout", 2*time.Second)

	// Cancel TimeOut
	// Cancel 能够立马生效
	c, _ = context.WithTimeout(context.Background(), time.Second)
	o = otherContext{c}
	c, cancel := context.WithTimeout(o, 2*time.Second)
	cancel()
	time.Sleep(100 * time.Millisecond) // let cancelation propagate
	select {
	case <-c.Done():
	default:
		fmt.Errorf("<-c.Done() blocked, but shouldn't have")
	}
	if e := c.Err(); e != context.Canceled {
		fmt.Errorf("c.Err() == %v want %v", e, context.Canceled)
	}

}

func testDeadline(c context.Context, name string, failAfter time.Duration) {
	select {
	case <-time.After(failAfter):
		fmt.Errorf("%s: context should have timed out", name)
	case <-c.Done():
		// context closed
	}
	if e := c.Err(); e != context.DeadlineExceeded {
		fmt.Errorf("%s: c.Err() == %v; want %v", name, e, context.DeadlineExceeded)
	}
}
