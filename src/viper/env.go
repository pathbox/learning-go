package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("spf")
	viper.BindEnv("id")

	os.Setenv("SPF_ID", 13)
	id := viper.Get("id")
	fmt.Println(id)
}
