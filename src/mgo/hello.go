package main

import (
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct{
	Name string
	Phone string
}

func main() {
	session, err := mgo.Dial("10.0.6.22")
	if err != nil{
		log.Fatal(err)
	}
	defer session.CLose()

	session.SetMode(mgo.Monotonic).C("people")

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 55 8989 7878"},
		&Person{"Cla", "+55 64 4342 5533"})
	if err != nil{
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}
