package main

import (
	"fmt"
)

func main() {
	saveUser(&user{Name: "Joe"})
}

type user struct {
	Name string
}

func (u *user) String() string {
	return u.Name
}

func saveUser(user *user) {
	withTx(func() {
		fmt.Println("保存用户: ", user.Name)
	})
}

func withTx(fn func()) {
	fmt.Println("Start...")
	fn()
	fmt.Println("End...")
}
