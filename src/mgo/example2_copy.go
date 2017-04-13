package main

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
	"net/http"
)

var (
	Database *mgo.Database
)

// listEntries lists all posts
func listPosts(w http.ResponseWriter, r *http.Request) {
	// We have a rather long running unit of work
	// (reading and listing all posts)
	// So it is worth copying the session
	collection := Database.C("posts").With(Database.Session.Copy())

	post := bson.D{}
	posts := collection.Find(bson.M{}).Iter()
	for posts.Next(&post) {
		// ...
	}
}

func main() {
	session, _ := mgo.Dial("mongodb://localhost:27017")
	Database := session.DB("myDb")

	// Count is a rather fast operation
	// No need to copy the session here
	count, _ := Database.C("posts").Count()
	fmt.Printf("Currently %d posts in the database", count)

	http.HandleFunc("/posts", listPosts)
	http.ListenAndServe(":8080", nil)
}
