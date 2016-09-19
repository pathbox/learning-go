package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()

  router.GET("/someGet", func(c *gin.Context){
    c.String(http.StatusOK, "Hello GET")
  })

  router.POST("/somePost", func(c *gin.Context){
    c.String(http.StatusOK, "Hello POST")
  })

  router.PUT("/somePut", func(c *gin.Context){
    c.String(http.StatusOK, "Hello PUT")
  })

  router.DELETE("/someDelete", func(c *gin.Context){
    c.String(http.StatusOK, "Hello DELETE")
  })

  router.PATCH("/somePatch", func(c *gin.Context){
    c.String(http.StatusOK, "Hello PATCH")
  })

  router.HEAD("/someHead", func(c *gin.Context){
    c.String(http.StatusOK, "Hello HEAD")
  })

  router.OPTIONS("/someOptions", func(c *gin.Context){
    c.String(http.StatusOK, "Hello OPTIONS")
  })

  // By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
