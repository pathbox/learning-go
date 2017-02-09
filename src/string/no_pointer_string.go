package main

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Environment string
	Version     string
	HostName    string
}

func (c *Config) String() string {
	return fmt.Sprintf("Environment: '%v'\nVersion:'%v'\nHostName: '%v'",
		c.Environment, c.Version, c.HostName)
}

func main() {

	jsonDoc := `
        {
            "Environment" : "Dev",
            "Version" : ""
        }`

	conf := &Config{}
	json.Unmarshal([]byte(jsonDoc), conf)
	fmt.Println(conf) // Prints
	//   Environment: 'Dev'
	//   Version:''
	//   HostName: ''

}
