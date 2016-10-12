package main

import (
	"fmt"
	"log"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

func init(){
  etcd.Register()
}

func main() {
  client := "localhost:2379"

  kv, err := libkv.NewStore(
    store.ETCD,
    []string{client},
    &store.Config{
      ConnectionTimeout: 10 * time.Second,
    },
  )
  if err != nil {
    log.Fatal("Cannot create store consul")
  }

  key := "foo"
  err = kv.Put(key, []byte("bar"), nil)
  if err != nil {
		_ = fmt.Errorf("Error trying to put value at key: %v", key)
	}

  pair, err := kv.Get(key)
  if err != nil {
		_ = fmt.Errorf("Error trying accessing value at key: %v", key)
	}
  fmt.Println("value: ", string(pair.Value))

	err = kv.Delete(key)
	if err != nil {
		_ = fmt.Errorf("Error trying to delete key %v", key)
	}

	fmt.Println("value: ", string(pair.Value))
}
