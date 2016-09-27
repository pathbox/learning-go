go client demo:

origin := "http://localhost:8080/"
url := "ws://localhost:8080/ws"
ws, err := websocket.Dial(url, "", origin)
if err != nil {
    log.Fatal(err)
}
if _, err := ws.Write([]byte("hello, world\n")); err != nil {
    log.Fatal(err)
}
var msg = make([]byte, 512)
var n int
if n, err = ws.Read(msg); err != nil {
    log.Fatal(err)
}
log.Printf("Received: %s.", msg[:n])

package main

import (
    "io"
    "log"
    "net/http"

    "github.com/ant0ine/go-json-rest/rest"
    "golang.org/x/net/websocket"
)

func main() {
  wsHandler := websocket.Handler(func(ws *websocket.Conn){
    io.Copy(ws, ws)
  })
  router, err := rest.MakeRouter(
        rest.Get("/ws", func(w rest.ResponseWriter, r *rest.Request) {
            wsHandler.ServeHTTP(w.(http.ResponseWriter), r.Request)
        }),
    )
  if err != nil{
    log.Fatal(err)
  }

  api := rest.NewApi()
  api.Use(rest.DefaultDevStack...)
  api.SetApp(router)
  log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))

}
