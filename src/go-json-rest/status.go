// Inspired by memcached "stats", this optional feature can be enabled to help monitoring the service.
// This example shows how to enable the stats, and how to setup the /.status route.
package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
)

func main() {
  api := rest.NewApi()
  statusMw := &rest.StatusMiddleware{}
  api.Use(statusMw)
  api.Use(rest.DefaultDevStack...)
  router, err := rest.MakeRouter(
    rest.Get("/.status", func(w rest.ResponseWriter, r *rest.Request){
      w.WriteJson(statusMw.GetStatus())
    }),
  )
  if err != nil {
    log.Fatal(err)
  }
  api.SetApp(router)
  log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
