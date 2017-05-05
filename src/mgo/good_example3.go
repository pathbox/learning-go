package main

import (
	"fmt"
	"github.com/pschlump/jsonp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
)

var MgoSession, _ = mgo.Dial("mongodb://127.0.0.1:27017")

type CompanyData struct {
	Id           bson.ObjectId "_id,omitempty"
	Company_Code string        // code
	Status       bool          // 开关, 关闭后不允许使用
	Max          int           // 最大在线人数 默认(不存在)为 500
}

type CustomerTargetData struct {
	Id           int
	Company_Code string
	Name         string
}

func main() {
	// MgoSession, _ := mgo.Dial("mongodb://127.0.0.1:27017")
	MgoSession.SetPoolLimit(400)
	log.Println("Here comming request")
	http.HandleFunc("/", getKey)
	http.HandleFunc("/create", createCustomerTarget)
	http.HandleFunc("/jsonp", jsonPfunc)
	log.Fatal(http.ListenAndServe(":9090", nil))

}

func getKey(w http.ResponseWriter, r *http.Request) {
	log.Println("Test mongo")
	data := &CompanyData{}
	bsonM := bson.M{"company_code": "aajd"}

	ok := Find(bsonM, data)
	log.Println("ok======: ", ok)

	log.Println("company: ", data.Id, "--", data.Company_Code)
	result := fmt.Sprintf("company id: %v -- company_code: %s", data.Id, data.Company_Code)

	w.Write([]byte(result))

}

func createCustomerTarget(w http.ResponseWriter, r *http.Request) {
	log.Println("Test create customer target")
	data := &CustomerTargetData{Id: 1, Name: "target_1", Company_Code: "aajd"}
	filter := bson.M{"id": data.Id}
	result := Upsert(filter, data)
	log.Println("create result: ", result)
}

func Find(filter, data interface{}) bool {
	session := MgoSession.Clone()
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	collection := session.DB("udesk_vistor_dev").C("company")
	err := collection.Find(filter).One(data)
	return err == nil
}

func Upsert(filter, data interface{}) bool {
	session := MgoSession.Clone()

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection := session.DB("udesk_vistor_dev").C("customer_target")
	_, _err := collection.Upsert(filter, data)
	log.Println("Upsert : %v [%v] %v", filter, _err, data)
	return _err == nil
}

func jsonPfunc(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, jsonp.JsonP("", rw, req))
}

//  lsof -i :27017 | wc -l  ＝> 490
// 但是　内存只消耗了70M　说明现在连接池中只有少数的连接在被使用，其他都是闲置连接
