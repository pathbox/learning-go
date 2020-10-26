package main

import "fmt"

type Token struct {
	Type int
	Raw  string
}

// 实现String方法 则可以被%s格式化输出
func (t Token) String() string {
	return fmt.Sprintf("%v<%q>", t.Type, t.Raw)
}

func main() {
	t := Token{
		Type: 1,
		Raw:  "good",
	}
	fmt.Printf("The token: %s", t)
}
