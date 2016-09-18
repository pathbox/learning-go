package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()
  router.POST("/from_post", func(c *gin.Context){
    message := c.PostForm("message")
    nick := c.DefaultPostForm("nick", "anonymous")
    c.JSON(200, gin.H{
      "status": "Posted",
      "message": message,
      "nick": nick,
    })
  })
  router.Run(":8080")
}
