import (
  "fmt"
  "github.com/levigross/grequests"
)

func main() {
  resp, err := grequests.Get("http://httpbin.org/get", nil)
  if err != nil {
    fmt.Println("Unable to make request: ", err)
  }
  fmt.Println(resp.String())
}