package main

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/DeanThompson/ginpprof"
)

func main() {
	router := gin.Default()
	runtime.GC()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.Wrap(router)

	// ginpprof also plays well with *gin.RouterGroup
	// group := router.Group("/debug/pprof")
	// ginpprof.WrapGroup(group)

	n := runtime.NumGoroutine()
	fmt.Println(n)
	router.Run(":8080")
}
