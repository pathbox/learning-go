/*
将[0:mid]的部分存储另一个数组tmp， 然后数组tmp出栈，和[mid:]之后的数据进行一一比对，如果一一比对相同，则说明是回文字符串
*/

package main

import "fmt"

func main() {
	str := []string{"a", "h", "a", "h", "a"}
	var tmp [100]string

	l := len(str)
	mid := l/2 - 1

	for i := 0; i <= mid; i++ {
		tmp[i] = str[i]
	}

	var next int

	if l%2 == 0 {
		next = mid + 1
	} else {
		next = mid + 2
	}

	top := mid

	for i := next; i < l-1; i++ {
		if str[i] != tmp[top] {
			break
		}
		top--
	}

	if top == 0 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
