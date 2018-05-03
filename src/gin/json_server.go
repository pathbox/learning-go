package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Binding from JSON
type Student struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Desc     Info   `json:"info" binding:"required"`
}

type Info struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/show", func(c *gin.Context) {
		var stu Student
		err := c.ShouldBindJSON(&stu)
		if err != nil {
			log.Println(err)
		}
		log.Println(stu)
		// if err := c.ShouldBindJSON(&json); err == nil {
		// 	if json.User == "manu" && json.Password == "123" {
		// 		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		// 	} else {
		// 		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		// 	}
		// } else {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// }

		c.JSON(http.StatusOK, gin.H{"status": "you are good"})
	})

	// Example for binding a HTML form (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err == nil {
			if form.User == "manu" && form.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
