package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()

  router.GET("/user/:name", func(c *gin.Context){
    name := c.Param("name")
    c.String(http.StatusOK, "Hello %s", name)
  })
  // However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
  router.GET("/user/:name/*action", func(c *gin.Context) {
    name := c.Param("name")
    action := c.Param("action")
    message := name + " is " + action
    c.String(http.StatusOK, message)
  })

  router.Run(":8080")
}
