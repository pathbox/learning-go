package main

import (
	"fmt"
	"html/template"
)

func main() {
	s := "Hello, <script>alert('you have been pwned')</script>!"

	// r1 := template.HTMLEscape(s)
	r2 := template.HTMLEscapeString(s)
	r3 := template.HTMLEscaper(s)
	fmt.Println("r1:", r2)
	fmt.Println("r2:", r3)
}

/*
对XSS最佳的防护应该结合以下两种方法：一是验证所有输入数据，有效检测攻击;另一个是对所有输出数据进行适当的处理，以防止任何已成功注入的脚本在浏览器端运行。

那么Go里面是怎么做这个有效防护的呢？Go的html/template里面带有下面几个函数可以帮你转义

func HTMLEscape(w io.Writer, b []byte) //把b进行转义之后写到w
func HTMLEscapeString(s string) string //转义s之后返回结果字符串
func HTMLEscaper(args ...interface{}) string //支持多个参数一起转义，返回结果字符串
*/
