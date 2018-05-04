package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// router.LoadHTMLGlob("/Users/pathbox/code/learning-go/src/gin/example/render_html/templates/*")
	router.LoadHTMLFiles("/Users/pathbox/code/learning-go/src/gin/example/render_html/templates/index.tmpl")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Index website",
		})
	})

	router.Run(":8080")
}
