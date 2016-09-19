package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.GET("/ping", func(c *gin.Context){
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  router.Run() // listen and server on 0.0.0.0:8080
}
