package main 

import (  
    "flag"
    "fmt"
    "os"
)

func main() {
  var help bool

  var helpText = "help text omitted for readability, shown in output instead."

  flag.BoolVar(&help, "help", false, helpText)

  flag.Parse()

  if help {
    fmt.Printf("%s\n", flag.Lookup("help").Usage)
    os.Exit(0)
  }
}