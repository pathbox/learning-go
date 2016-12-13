package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	input := `
  {
    "Mysql":{
      "Host":"127.0.0.1",
      "Port":"3306",
      "User":"root",
      "Passwd":"",
      "Db":"mob_portal"
    },
    "HashKey":"lsidieieisi"
  }
  `
	var info struct {
		Mysql struct {
			Host, Port, User, Passwd, Db string
		}
		HashKey string
	}
	err := json.Unmarshal([]byte(input), &info)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", info.Mysql.Host)
}
