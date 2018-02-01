package main

import "fmt"

// func main() {
// 	a := []int{1, 2, 3}
// 	for _, i := range a {
// 		fmt.Println(&i)
// 		func() {
// 			fmt.Println(&i)
// 		}()
// 	}
// }

// 这个就是闭包的“神奇”之处。闭包里的非传递参数外部变量值是传引用的，
// 在闭包函数里那个i就是外部非闭包函数自己的参数，所以是相当于引用了外部的变量，
//  i 的值执行到第三次是3 ，闭包是地址引用所以打印了3次i地址指向的值，所以是3，3，3

func main() {
	a := []int{1, 2, 3}
	for _, i := range a {
		fmt.Println(&i)

		fput(i)()
	}
}

func fput(i int) func() {
	return func() {
		fmt.Println(&i)
	}
}
