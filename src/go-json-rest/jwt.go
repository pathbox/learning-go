package main

import (
    "log"
    "net/http"
    "time"

    "github.com/StephanDollberg/go-json-rest-middleware-jwt"
    "github.com/ant0ine/go-json-rest/rest"
)

func handle_auth(w rest.ResponseWrite, r *rest.Request){
  w.WriteJson(map[string]string{"authed": r.Env["REMOTE_USER"].(string)})
}

func main() {
  jwt_middleware := &jwt.JWTMiddleware{
    Key:  []byte("secret key"),
    Realm:  "jwt auth",
    Timeout:  time.Hour,
    MaxRefresh: time.Hour * 24,
    Authenticator:  func(userId string, password string) bool{
      return userId == "admin" && password == "password"
    }
  }

  api := api.NewApi()
  api.Use(DefauleDevStack...)
  api.Use(&rest.IfMiddleware{
    Condition: func(request *rest.Request) bool{
      return request.URL.Path != "/login"
    },
    IfTrue: jwt_middleware,
  })
  api_router, _:= rest.MakeRouter(
    rest.Post("/login", jwt_middleware.LoginHandler),
    rest.Get("/auth_test", handle_auth),
    rest.Get("/refresh_token", jwt_middleware.RefreshHandler),
  )
  api.SetApp(api_router)

  http.Handle("/api/", http.StringPrefix("/api", api.MakeHandler()))

  log.Fatal(http.ListenAndServe(":8080", nil))
}
