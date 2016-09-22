package main

import (
  "log"
  "net/http"
  "github.com/ant0ine/go-json-rest/rest"
)

func main() {
  api := rest.NewApi()
  api.Use(rest.DefaultDevStack...)
  api.SetApp(rest.AppSimple(func(w rest.ResponseWriter, r *rest.Request){
    w.WriteJson(map[string]string{"Body": "Hello World!"})
    }))
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
