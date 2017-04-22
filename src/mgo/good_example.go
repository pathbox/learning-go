package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type MongoConn struct {
	Session *mgo.Session
}

func (m *MongoConn) Connect() *mgo.Session {
	var err error
	if m.Session == nil {
		m.Session, err = mgo.Dial("mongodb://localhost:27017")
		if err != nil {
			panic(err)
		}
		m.Session.SetMode(mgo.Monotonic, true)
	}
	return m.Session.Clone()
}

var GlobalMgoSession *MongoConn

func Insert(db, collection string, records []bson.M) {
	s := GlobalMgoSession.Connect()
	defer s.Close()
	c := s.DB(db).C(collection)
	for _, rec := range records {
		if err := c.Insert(&rec); err != nil {
			panic(err)
		}
	}
}

type Results bson.M

func Get(db, collection string, spec bson.M, idx, limit int) []Results {
	var out []Results
	s := GlobalMgoSession.Connect()
	defer s.Close()
	c := s.DB(db).C(collection)
	if limit > 0 {
		err := c.Find(spec).Skip(idx).Limit(limit).All(&out)
	} else {
		err := c.Find(spec).Skip(idx).All(&out)
	}
	if err != nil {
		panic(err)
	}
	return out
}

// helper function to query MongoDB (test.collection database).
func processRequest(query string) map[string]interface{} {
	response := make(map[string]interface{})
	var spec bson.M
	data := Get("test", "collection", spec, 0, -1)
	response["data"] = data
	response["query"] = query
	data = nil
	return response
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	response := processRequest(query)
	js, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	js = nil
	response = nil
}

func Server() {
	http.HandleFunc("/", RequestHandler)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	Server()
}
