package main

import (
	"fmt"
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
		m.Session.SetPoolLimit(400) // 设置　连接池上限　实际测试其实没有卵用,该超过还是会超过
		m.Session.SetMode(mgo.Monotonic, true)
	}
	return m.Session.Clone()
}

var GlobalMgoSession *MongoConn

type CompanyData struct {
	Id           bson.ObjectId "_id,omitempty"
	Company_Code string        // code
	Status       bool          // 开关, 关闭后不允许使用
	Max          int           // 最大在线人数 默认(不存在)为 500
}

func main() {
	GlobalMgoSession = &MongoConn{}
	log.Println("Here comming request")
	http.HandleFunc("/", getKey)
	log.Fatal(http.ListenAndServe(":9090", nil))

}

func getKey(w http.ResponseWriter, r *http.Request) {
	log.Println("Test mongo")
	data := &CompanyData{}
	bsonM := bson.M{"company_code": "aajd"}

	ok := GlobalMgoSession.Find(bsonM, data)
	log.Println("ok======: ", ok)

	log.Println("company: ", data.Id, "--", data.Company_Code)
	result := fmt.Sprintf("company id: %v -- company_code: %s", data.Id, data.Company_Code)

	w.Write([]byte(result))

}

func (m *MongoConn) Find(filter, data interface{}) bool {
	session := m.Connect()
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	collection := session.DB("udesk_vistor_dev").C("company")
	err := collection.Find(filter).One(data)
	return err == nil
}
