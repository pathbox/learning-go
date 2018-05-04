package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
)

// 实际上就是自己实现 gin 的middleware， 在middleware中 进行txn defer txn.End()
func NewrelicMiddleware(appName, key string) gin.HandlerFunc {
	if appName == "" || key == "" {
		return func(c *gin.Context) {}
	}

	config := newrelic.NewConfig(appName, key)
	app, err := newrelic.NewApplication(config)

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		txn := app.StartTransaction(c.Request.URL.Path, c.Writer, c.Request)
		defer txn.End()
		c.Next()
	}
}

func main() {
	router := gin.Default()

	router.Use(NewrelicMiddleware("APP-NAME", "APP_KEY"))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It Works")
	})

	router.Run()
}
