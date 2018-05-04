package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("form_post", FormPost)
	router.POST("form_post1", FormPost1)

	router.Run(":8080")
}

func FormPost(c *gin.Context) {
	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
	})
}

func FormPost1(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name")
	message := c.PostForm("message")

	fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
}
