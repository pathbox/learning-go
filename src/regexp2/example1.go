package main

import (
	"fmt"

	"github.com/dlclark/regexp2"
)

// 示例
func main() {
	text := `Hello 世界！123 Go.`

	// 查找连续的小写字母
	reg := regexp2.MustCompile(`[a-z]+`, 2)
	m, _ := reg.FindStringMatch(text)
	// you can get all the groups too
	for _, gp := range m.Groups() {
		fmt.Println("Group:", gp.Capture.String())
	}

	// // 查找连续的非小写字母
	matches := []string{}
	reg = regexp2.MustCompile(`[^a-z]+`, 110)
	m, _ = reg.FindStringMatch(text)
	// you can get all the groups too
	for m != nil {
		matches = append(matches, m.String())
		m, _ = reg.FindNextMatch(m)
	}
	fmt.Println(matches)

	// // ["H" " 世界！123 G" "."]

	// // 查找连续的单词字母
	// reg = regexp2.MustCompile(`[\w]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello" "123" "Go"]

	// // 查找连续的非单词字母、非空白字符
	// reg = regexp2.MustCompile(`[^\w\s]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["世界！" "."]

	// // 查找连续的大写字母
	// reg = regexp2.MustCompile(`[[:upper:]]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["H" "G"]

	// // 查找连续的非 ASCII 字符
	// reg = regexp2.MustCompile(`[[:^ascii:]]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["世界！"]

	// // 查找连续的标点符号
	// reg = regexp2.MustCompile(`[\pP]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["！" "."]

	// // 查找连续的非标点符号字符
	// reg = regexp2.MustCompile(`[\PP]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello 世界" "123 Go"]

	// // 查找连续的汉字
	// reg = regexp2.MustCompile(`[\p{Han}]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["世界"]

	// // 查找连续的非汉字字符
	// reg = regexp2.MustCompile(`[\P{Han}]+`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello " "！123 Go."]

	// // 查找 Hello 或 Go
	// reg = regexp2.MustCompile(`Hello|Go`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello" "Go"]

	// // 查找行首以 H 开头，以空格结尾的字符串
	// reg = regexp2.MustCompile(`^H.*\s`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello 世界！123 "]

	// // 查找行首以 H 开头，以空白结尾的字符串（非贪婪模式）
	// reg = regexp2.MustCompile(`(?U)^H.*\s`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello "]

	// // 查找以 hello 开头（忽略大小写），以 Go 结尾的字符串
	// reg = regexp2.MustCompile(`(?i:^hello).*Go`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello 世界！123 Go"]

	// // 查找 Go.
	// reg = regexp2.MustCompile(`\QGo.\E`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Go."]

	// // 查找从行首开始，以空格结尾的字符串（非贪婪模式）
	// reg = regexp2.MustCompile(`(?U)^.* `)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello "]

	// // 查找以空格开头，到行尾结束，中间不包含空格字符串
	// reg = regexp2.MustCompile(` [^ ]*$`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // [" Go."]

	// // 查找“单词边界”之间的字符串
	// reg = regexp2.MustCompile(`(?U)\b.+\b`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello" " 世界！" "123" " " "Go"]

	// // 查找连续 1 次到 4 次的非空格字符，并以 o 结尾的字符串
	// reg = regexp2.MustCompile(`[^ ]{1,4}o`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello" "Go"]

	// // 查找 Hello 或 Go
	// reg = regexp2.MustCompile(`(?:Hell|G)o`)
	// fmt.Printf("%q\n", reg.FindStringMatch(text, -1))
	// // ["Hello" "Go"]

	// // 查找 Hello 或 Go，替换为 Hellooo、Gooo
	// reg = regexp2.MustCompile(`(?PHell|G)o`)
	// fmt.Printf("%q\n", reg.ReplaceAllString(text, "${n}ooo"))
	// // "Hellooo 世界！123 Gooo."

	matches = []string{}
	t := "oss:dsds:GetBucketStat"
	reg = regexp2.MustCompile(`oss:.*:GetBucketStat`, 0)
	isMatch, _ := reg.MatchString(t)
	fmt.Println(isMatch)
	// // 特殊字符的查找
	// reg = regexp2.MustCompile(`[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\$\.\*\+\?\{\}\(\)\[\]\|]`)
	// fmt.Printf("%q\n", reg.ReplaceAllString("\f\t\n\r\v\123\x7F\U0010FFFF\\^$.*+?{}()[]|", "-"))
	// // "----------------------"
}
