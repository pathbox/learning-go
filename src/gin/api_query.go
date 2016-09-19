package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func main() {
  // Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe

  router := gin.Default()
  router.GET("/welcome", func(c *gin.Context){
    firstname := c.DefaultQuery("firstname", "Guest")
    lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
    c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
  })

  router.Run(":8080")
}
