package main

import (
	"github.com/bitly/go-simplejson"
)

func main() {
	contents, _ := ioutil.ReadAll(res.Body)
	json, err := simplejson.NewJson(contents)
	var nodes = make(map[string]interface{})
	nodes, _ = json.Map()
	// 如果map的value仍是一个json机构，回到第二步
	for key, _ := range nodes {

	}
}
