package main

import (
	"fmt"
	"os"
)

func main() {
	gm := os.Getenv("GO111MODULE")
	fmt.Println(gm)

	str, err := os.LookupEnv("GO111MODULE")
	if !err {
		panic(err)
	}
	fmt.Println(str)

	_ = os.Unsetenv("DB_XXXCCC")
	_ = os.Setenv("DB_XXXCCC", "db:/user3@example")
	str3 := os.Getenv("DB_XXXCCC")
	fmt.Println(str3)

	os.Setenv("GO111MODULE", "off")
	gm1 := os.Getenv("GO111MODULE")
	fmt.Println(gm1)
}
