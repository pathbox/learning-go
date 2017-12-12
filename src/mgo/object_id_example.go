package main

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type Work struct {
	Id bson.ObjectId `json:"-" bson:"_id"`
	IdHex string
	Name string
}

func main(){
	work := &Work{}
	work.IdHex = bson.NewObjectId().Hex()
	fmt.Println("object id", work.IdHex)
	work.Id = bson.ObjectIdHex(work.IdHex)
	work.Name = "Hello"

	fmt.Println("work struct: ", work)
}