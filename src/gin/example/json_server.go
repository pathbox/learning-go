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
	Data     Info   `json:"info" binding:"required"`
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
		// err := c.ShouldBindJSON(&stu)
		err := c.BindJSON(&stu)
		if err != nil {
			log.Println("Error:", err)
		}
		log.Println(stu)

		c.JSON(http.StatusOK, stu)
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":9090")
}
