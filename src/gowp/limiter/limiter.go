package limiter

import (
	"time"

	"github.com/xxjwxc/public/myredis"
)

type mySemaphore struct {
	limit       int
	nameSpace   string
	isTsTimeout bool
	redisClient myredis.RedisDial
}

type LimiterIFS interface {
	Acquire(timeout int) (string, error)
	Release(token string)
	GetTimeDuration(token string) (time.Duration, error)
}

func NewLimiter(ops ...Option) (lifs LimiterIFS) {
	var tmp = mySemaphore{}
	for _, o := range ops {
		o.apply(&tmp)
	}

	// set default
	if tmp.limit <= 0 {
		tmp.limit = 1
	}
	if len(tmp.nameSpace) == 0 {
		tmp.nameSpace = "gowp/nameSpace"
	}

	//-------------end

	if tmp.redisClient == nil { // cache模式
		lifs = &limiterCache{mySemaphore: tmp}
	} else { // redis sync
		lifs = &limiterRedis{mySemaphore: tmp}
	}
	lifs.Init()

	return
}
