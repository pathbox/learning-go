package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	s1 := `{"result":"","status":false,"code":"1028","msg":"字段"}`
	s2 := `{"result":"","status":false,"code":"1028","msg":"字段无效"}`
	var m1, m2 interface{}
	json.Unmarshal([]byte(s1), &m1)
	json.Unmarshal([]byte(s2), &m2)
	m := make(map[string]interface{})
	m["comMsg"] = []interface{}{m1, m2}
	s, _ := json.Marshal(m)
	fmt.Printf("%s", s)
}
