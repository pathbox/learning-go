package main

import "fmt"

var PM map[string]*Person

type Person struct {
	Age  int
	Name string
}

func main() {
	PM = make(map[string]*Person) // Don't forget it
	PM["Joe"] = &Person{Age: 28, Name: "Joe"}
	DoAction()
}

func DoAction() {
	j := PM["Joe"]            // 先将值读出赋值给j
	delete(PM, "Joe")         // 在map中delete Joe
	fmt.Printf("Joe:%v\n", j) // 打印输出 j 不报错
	fmt.Printf("PM:%v\n", PM)
	fmt.Println(j.Name)
}
