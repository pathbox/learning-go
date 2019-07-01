package main

import(
	"fmt"
	"sync"
)

type User struct {
	Name string
}

func main() {
	user := &User{
		Name: "Cary",
	}
	var wg sync.WaitGroup
	var wg1 sync.WaitGroup

	go func(u *User) {
		wg.Add(1)
		u.Name = "New Guy" // u is just a copy from user
		fmt.Println("My name is:",u.Name)
	}(user)
	wg.Wait()
	fmt.Println("Now My name is:",user.Name)

	go func() {
		wg1.Add(1)
		user.Name = "New Guy" // u is just a copy from user
		fmt.Println("My haha name is:",user.Name)
	}()
	wg1.Wait()
	fmt.Println("Now haha My name is:",user.Name)

}

// 通过闭包传进的变量是copy拷贝，直接在闭包中使用的变量不是copy,如果在闭包中进行修改，怎修改的是原地址的值，而不是copy的值