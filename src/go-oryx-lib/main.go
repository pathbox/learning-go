package main

import(
  "log"
  "net/http"
  ohttp "github.com/ossrs/go-oryx-lib/http"
)

type Patload struct {
  Foo string
  Bar string
}

const (
  errorSystemError ohttp.SystemError = 100 + iota
  errorSystemComplexError
)

func main() {
  ohttp.Server = "akserver"

  data := "data string"

  http.HandleFunc("/", ohttp.Data(nil, data).ServeHTTP)

  data1 := Payload{
    Foo: "foo",
    Bar: "bar",
  }
  http.HandleFunc("/data_struct", ohttp.Data(nil, data1).ServHTTP)

  http.HandleFunc("/system_error", ohttp.Error(nil, errorSystemError).ServeHTTP)

  sce := ohttp.SystemComplexError{
    Code: errorSystemComplexError,
    Message: "SystemComplexError string",
  }
  http.HandleFunc("/system_complex_error", ohttp.Error(nil, sce).ServeHTTP)

  err := http.ListenAndServe(":9090", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
