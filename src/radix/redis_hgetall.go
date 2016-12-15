package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/mediocregopher/radix.v2/redis"
)

type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reply, err := conn.Cmd("HGETALL", "album:1").Map()
	if err != nil {
		log.Fatal(err)
	}

	ab, err := populateAlbum(reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ab)
}

func populateAlbum(reply map[string]string) (*Album, error) {
	var err error
	ab := new(Album)
	ab.Title = reply["title"]
	ab.Artist = reply["artist"]
	ab.Price, err = strconv.ParseFloat(reply["price"], 64)
	if err != nil {
		return nil, err
	}
	ab.Likes, err = strconv.Atoi(reply["likes"])
	if err != nil {
		return nil, err
	}
	return ab, nil
}
