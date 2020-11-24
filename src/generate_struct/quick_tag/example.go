package main

import (
    "encoding/json"
    "fmt"
    . "github.com/dovejb/quicktag"
    "reflect"
)

type Person struct {
    Name       string
    Age        int
    MyChildren []Person
}

func main() {
    p := Person{
        Name: "dovejb",
        Age:  6,
        MyChildren: []Person{
            Person{
                Name: "baby",
                Age:  3,
            },
        },
    }

    var p2 Person

    buf, _ := json.Marshal(Q(p))
    fmt.Println(string(buf))
    // {"name":"dovejb","age":6,"my_children":[{"name":"baby","age":3,"my_children":null}]}

    json.Unmarshal(buf, Q(&p2))
    fmt.Println(reflect.DeepEqual(p, p2))
    // true
}