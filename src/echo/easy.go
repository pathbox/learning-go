package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

func main() {
	e := echo.New()
	e.Get("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.POST("/users", saveUser)
	e.Get("/users/:id", getUser)
	e.PUT("/users/:id", putUser)
	e.DELETE("/users/:id", deleteUser)
	e.Run(standard.New(":9090"))
}

func show(c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")
}

func save(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
}

func save(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")

	avatar, err := c.FormFile("avatar")
	if err != nil {
		panic(err)
	}

	src, err := avatar.Open()
	if err != nil {
		return err
	}

	defer src.Close()
	dst, err := os.Create(avatar.Filename)
	if err != nil {
		panic(err)
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "<b> Thank you </b>")
}

type User struct {
	Name string
	Email string
}

e.POST("/users", func(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {  // Bind JSON or XML or form payload into Go struct based on Content-Type request header.
		return err
	}
	return c.JSON(http.StatusCreated, u) // Render response as JSON or XML with status code.
})

// Server any file from static directory for path /static/*.
e.Static("/static", "static")

// Middleware

// Root level middleware
e.Use(middleware.Logger())
e.Use(middleware.Recover())

// Group level middleware

g := e.Group("/admin")
g.Use(middleware.BasicAuth(func(username, password string) bool {
	if username == "admin" && password == "password" {
		return true
	}
	return false
}))

// Route level middleware

track := func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("request to /users")
		return next()
	}
}

e.GET("/users", func(c echo.Context) error {
	return c.String(http.StatusOK, "/users")
	}, track)








































