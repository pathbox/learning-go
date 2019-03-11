package gocachebenchmarks

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/allegro/bigcache"
	"github.com/bluele/gcache"
	"github.com/coocood/freecache"
	cache2go "github.com/muesli/cache2go"
	cache "github.com/patrickmn/go-cache"
)

func BenchmarkCache2Go(b *testing.B) {
	c := cache2go.Cache("test")

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("%20d", i)
			c.Add(fmt.Sprintf("item%d", i), 1*time.Minute, value)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, err := c.Value(fmt.Sprintf("item%d", i))
			if err == nil {
				_ = value
			}
		}
	})
}

func BenchmarkGoCache(b *testing.B) {
	c := cache.New(1*time.Minute, 5*time.Minute)

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("%20d", i)
			c.Add(fmt.Sprintf("item%d", i), value, cache.DefaultExpiration)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, found := c.Get(fmt.Sprintf("item%d", i))
			if found {
				_ = value
			}
		}
	})
}

func BenchmarkFreecache(b *testing.B) {
	c := freecache.NewCache(1024 * 1024 * 5)

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("%20d", i)
			c.Set([]byte(fmt.Sprintf("item%d", i)), []byte(value), 60)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, err := c.Get([]byte(fmt.Sprintf("item%d", i)))
			if err == nil {
				_ = value
			}
		}
	})
}

func BenchmarkBigCache(b *testing.B) {
	c, _ := bigcache.NewBigCache(bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,
		// time after which entry can be evicted
		LifeWindow: 10 * time.Minute,
		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,
		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,
		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 10,
	})

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("%20d", i)
			c.Set(fmt.Sprintf("item%d", i), []byte(value))
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, err := c.Get(fmt.Sprintf("item%d", i))
			if err == nil {
				_ = value
			}
		}
	})
}

func BenchmarkGCache(b *testing.B) {
	c := gcache.New(b.N).LRU().Build()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("%20d", i)
			c.Set(fmt.Sprintf("item%d", i), value)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, err := c.Get(fmt.Sprintf("item%d", i))
			if err == nil {
				_ = value
			}
		}
	})
}

// No expire, but helps us compare performance
func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value := fmt.Sprintf("%20d", i)
			m.Store(fmt.Sprintf("item%d", i), value)
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			value, found := m.Load(fmt.Sprintf("item%d", i))
			if found {
				_ = value
			}
		}
	})
}
