// 首先，我们有一个固定容量的桶，有水进来，也有水出去。对于流进来的水，我们无法预计共有多少水流进来，也无法预计流水速度，但
// 对于流出去的水来说，这个桶可以固定水流的速率，而且当桶满的时候，多余的水会溢出来

package main

import (
  "fmt"
  "time"
)

func main() {

  timeUnix := time.Now().Unix()
  water := 0
  rate := 0.0
  capacity := 10000

  rate_bucket(timeUnix, rate, water, capacity)

}

func rate_bucket(timeUnix int64, rate float32, water int32, capacity int32){
  now = time.Now().Unix()

  tmp_water := water - ( now - timeUnix) * rate

  timeUnix = now
  if tmp_water+1 > 0 {

    water = tmp_water + 1
    return true
    } else {
    return false
    }
}

// 限制原理是 水桶加满速度和漏水速度的比较， 如果水桶加满了，则限制产生