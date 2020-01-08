package ratelimiter

//Bucket算法接口
type RateLimiter interface {
	Wait(n int64)
	Rate() float32
	Capacity() int64
}

type RateLimit struct {
	strategy    string
	ratelimiter RateLimiter
}

func NewRateLimit(strategy string, rate float32, capacity int64) *RateLimit {
	if rate <= 0.0 || capacity <= 0 {
		return nil
	}
	var ratelimiter RateLimiter
	if strategy == "leaky" {
		ratelimiter = NewLeakyBucketWithRate(rate, capacity)
	} else if strategy == "token" {
		ratelimiter = NewBucketWithRate(rate, capacity)
	} else {
		return nil
	}
	return &RateLimit{
		strategy:    strategy,
		ratelimiter: ratelimiter,
	}
}

func (rl *RateLimit) Stop(n int64) {
	rl.ratelimiter.Wait(n)
	return
}

func (rl *RateLimit) Strategy() string {
	return rl.strategy
}

func (rl *RateLimit) Rate() float32 {
	return rl.ratelimiter.Rate()
}

func (rl *RateLimit) Capacity() int64 {
	return rl.ratelimiter.Capacity()
}
