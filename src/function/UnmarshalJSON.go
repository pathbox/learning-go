package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var jsonDoc = []byte(`["add", "sub", "mul", "div"]`)

var registry = map[string]binFunc{
	"add": func(x, y int) int { return x + y },
	"sub": func(x, y int) int { return x - y },
	"mul": func(x, y int) int { return x * y },
	"div": func(x, y int) int { return x / y },
}

type binFunc func(int, int) int

//实现了下面这个方法就可以自己反序列化自己
func (fn *binFunc) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	found := registry[name]
	if found == nil {
		return fmt.Errorf("unknow function in (*binFunc)UnmarshalJSON: %s", name)
	}

	*fn = found
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())
	var fns []binFunc
	if err := json.Unmarshal(jsonDoc, &fns); err != nil {
		log.Fatal(err)
	}

	fn := fns[rand.Intn(len(fns))]
	x := fn(12, 5)
	fmt.Println(x)
}
