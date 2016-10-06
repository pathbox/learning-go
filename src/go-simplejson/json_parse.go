package main

import(
  "fmt"
  "os"
  "github.com/bitly/go-simplejson"
)

func errHandle(err error)  {
  if err != nil{
    fmt.Println("error: ", err)
    os.Exit(1)
  }
}

func main() {
  js, err := simplejson.NewJson([]byte(`{
    "test":{
    "array": [1, "2", 3],
        "int": 10,
        "float": 5.150,
        "bignum": 9223372036854775807,
        "string": "simplejson",
        "bool": true
  }
    }`))

    errHandle(err)

    arr, _ := js.Get("test").Get("array").Array()
    i, _ := js.Get("test").Get("int").Int()
    ms := js.Get("test").Get("string").MustString()
    fmt.Println("array: ", arr)
	  fmt.Println("int: ", i)
	  fmt.Println("string: ", ms)
}
