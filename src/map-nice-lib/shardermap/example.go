package main

import (
	"fmt"

	"github.com/antoniomo/shardedmap"
)

type Data struct {
	ID string
	V  int
}

func main() {
	strmap := shardedmap.NewStrMap(64)
	a := Data{ID: "a", V: 1}
	b := Data{ID: "b", V: 2}

	strmap.Store(a.ID, a)
	strmap.Store(b.ID, b)

	a2, ok := strmap.Load(a.ID)
	if !ok {
		panic("ARGH!")
	}
	b2, ok := strmap.Load(b.ID)
	if !ok {
		panic("ARGH!")
	}

	fmt.Println("Same with range (note the random order):")
	strmap.Range(func(key string, value interface{}) bool {
		fmt.Printf("Key: %s, Value: %+v\n", key, value)
		return true
	})

	fmt.Println("LoadOrStore over a.ID")
	los, ok := strmap.LoadOrStore(a.ID, b)
	if !ok {
		panic("ARGH!")
	}
	fmt.Printf("Key: %s, Value: %+v\n", a.ID, los)

	strmap.Delete(a.ID)
	strmap.Delete(b.ID)
	_, ok = strmap.Load(a.ID)
	if ok {
		panic("ARGH!")
	}
	_, ok = strmap.Load(b.ID)
	if ok {
		panic("ARGH!")
	}
}