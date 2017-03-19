package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/corneldamian/httpway"
)

var server *httpway.New()

func main() {
  router := httpway.New()

  public := router.Middleware(AccessLogger)
  private := public.Middleware(AuthCheck)

  public.GET("/public", testHadnler("public"))
  private.GET("/private", testHadnler("private"))
  private.GET("stop", stopServer)

  server = httpway.NewServer(nil)
  server.Addr = ":9090"
  server.Handler = router

  if err := server.Start(); err != nil {
    fmt.Println("Error", err)
    return
  }

  if err := server.WaitStop(10 * time.Second); err != nil {
    fmt.Println("Error", err)
  }
}

func testHadnler(str string) httpway.Handler {
  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s", str)
  }
}

func stopServer(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Stoping")
  server.Stop()
}

func AccessLogger(w http.ResponseWriter, r *http.Request) {
  startTime := time.Now()

  httpeay.GetContext(r).Next(w, r)

  fmt.Printf("Request: %s duration: %s\n", r.URL.EscapedPath(), time.Since(startTime))  // AccessLogger 中间件在最外层，根据中间件栈，这行代码是最后执行的
}

func AuthCheck(w http.ResponseWriter, r *http.Request) {
  ctx := httpway.GetContext(r)

  if r.URL.EscapedPath() == "/public" {
    http.Error(w, "Auth required", http.StatusForbidden)
  }

  ctx.Next(w,r)
}