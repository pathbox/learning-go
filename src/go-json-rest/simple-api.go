package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
)

func main() {
  api := rest.NewApi()
  api.Use(rest.DefaultSDevStack...)

  router, err := rest.MakeRouter(
    rest.Get("/message", func(w rest.ResponseWriter, req *rest.Request){
      w.WtireJson(map[string]string{"Body": "Hello World!"})
    }),
  )

  if err != nil{
    log.Fatal(err)
  }
  api.SetApp(router)

  http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

  http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("."))))

  log.Fatal(http.ListenAndServe(":8080", nil))
}
//
// Combine Go-Json-Rest with other handlers.
//
// api.MakeHandler() is a valid http.Handler, and can be combined with other handlers. In this example the api handler is used under the /api/ prefix, while a FileServer is instantiated under the /static/ prefix.
//
// curl demo:
//
// curl -i http://127.0.0.1:8080/api/message
// curl -i http://127.0.0.1:8080/static/main.go
