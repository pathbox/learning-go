package main

import (
    "github.com/gin-gonic/gin"
)

type myForm struct {
    Colors []string `form:"colors[]"`
}

func main() {
    r := gin.Default()

    r.LoadHTMLGlob("views/*")
    r.GET("/", indexHandler)
    r.POST("/", formHandler)

    r.Run(":8080")
}

func indexHandler(c *gin.Context) {
    c.HTML(200, "form.html", nil)
}

func formHandler(c *gin.Context) {
    var fakeForm myForm
    c.Bind(&fakeForm)
    c.JSON(200, gin.H{"color": fakeForm.Colors})
}

<form action="/" method="POST">
    <p>Check some colors</p>
    <label for="red">Red</label>
    <input type="checkbox" name="colors[]" value="red" id="red" />
    <label for="green">Green</label>
    <input type="checkbox" name="colors[]" value="green" id="green" />
    <label for="blue">Blue</label>
    <input type="checkbox" name="colors[]" value="blue" id="blue" />
    <input type="submit" />
</form>