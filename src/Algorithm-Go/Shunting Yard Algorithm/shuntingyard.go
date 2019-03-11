package main

import (
	"fmt"
	"strings"
)

type Node struct {
	isNumber bool
	val      float64
	opt      byte
}

//调度场算法，返回后缀表达式即逆波兰式。
func ShuntingYard(str string) []*Node {
	//去除所有的空格。
	str = strings.Replace(str, " ", "", -1)
	//当前位是否可以输入数字。
	canNumber := true
	//当前数字的符号，因为有可能一上来就是一个-
	sign := 1
	//输出队列
	var queue []*Node
	//符号栈
	var stack []byte
	//读取每次的输入指针头。
	p := 0
	N := len(str)
	for p < N {
		//按字符读入的是数字或者小数点.
		if (str[p] >= '0' && str[p] <= '9') || str[p] == '.' {
			if canNumber {
				//当前正在读取的这个数的值
				var val float64 = 0
				//小数代表的权重
				var w float64 = 1
				//当前是否已经输入了小数点
				hasDot := false
				for p < N && ((str[p] >= '0' && str[p] <= '9') || str[p] == '.') {
					if str[p] == '.' {
						if hasDot {
							panic("一个数字不能有两个小数点")
						}
						hasDot = true
					} else {
						if hasDot {
							w *= 0.1
							val += float64(str[p]-'0') * w
						} else {
							val = val*10 + float64(str[p]-'0')
						}
					}
					p++
				}
				p--
				if str[p] == '.' {
					panic("一个小数点，不能作为数字")
				}
				queue = append(queue, &Node{isNumber: true, val: val * float64(sign)})
				sign = 1
				canNumber = false
			} else {
				panic(fmt.Sprintf("在%c附近表达有误", str[p]))
			}
		} else {
			switch str[p] {
			case '(':
				//入栈
				stack = append(stack, str[p])
			case ')':
				//出栈直到'('
				for len(stack) != 0 && stack[len(stack)-1] != '(' {
					//出栈
					queue = append(queue, &Node{isNumber: false, opt: stack[len(stack)-1]})
					stack = stack[:len(stack)-1]
				}
				//如果已经变空了，那么括号不匹配
				if len(stack) == 0 {
					panic("左右括号不匹配")
				}
				//将左括号丢弃。
				stack = stack[:len(stack)-1]
			case '+':
				fallthrough
			case '-':
				//如果这时候还能放数字。
				if canNumber {
					if str[p] == '-' {
						sign = sign * -1
					}
				} else {
					//+优先级低，因此栈顶的出栈，自己进栈。
					for len(stack) != 0 && (stack[len(stack)-1] == '*' || stack[len(stack)-1] == '/' || stack[len(stack)-1] == '+' || stack[len(stack)-1] == '-') {
						//出栈
						queue = append(queue, &Node{isNumber: false, opt: stack[len(stack)-1]})
						stack = stack[:len(stack)-1]
					}
					stack = append(stack, str[p])
					canNumber = true
				}
			case '*':
				fallthrough
			case '/':
				//如果这时候还能放数字。
				if canNumber {
					panic(fmt.Sprintf("在%c附近表达有误。", str[p]))
				} else {
					//如果同等优先级，那么先计算前面的。
					for len(stack) != 0 && (stack[len(stack)-1] == '*' || stack[len(stack)-1] == '/') {
						//出栈
						queue = append(queue, &Node{isNumber: false, opt: stack[len(stack)-1]})
						stack = stack[:len(stack)-1]
					}
					//优先级很高，直接入栈
					stack = append(stack, str[p])
					canNumber = true
				}
			default:
				panic(fmt.Sprintf("未能识别的符号：%c", str[p]))
			}
		}
		p++
	}
	//符号栈中的统统进去
	for len(stack) != 0 {
		if stack[len(stack)-1] == '(' {
			panic("左右括号不匹配。")
		}
		//出栈
		queue = append(queue, &Node{isNumber: false, opt: stack[len(stack)-1]})
		stack = stack[:len(stack)-1]
	}
	return queue
}

//打印逆波兰式
func PrintNodes(queue []*Node) {
	for i := 0; i < len(queue); i++ {
		if queue[i].isNumber {
			fmt.Printf("%v ", queue[i].val)
		} else {
			fmt.Printf("%c ", queue[i].opt)
		}
	}
	fmt.Println()
}

//计算逆波兰式
func CalNodes(queue []*Node) float64 {
	var stack []float64
	PrintNodes(queue)

	for _, v := range queue {
		if v.isNumber {
			//数字直接入栈
			stack = append(stack, v.val)
		} else {
			switch v.opt {
			case '+':
				tmp := stack[len(stack)-1] + stack[len(stack)-2]
				stack[len(stack)-2] = tmp
				stack = stack[:len(stack)-1]
			case '-':
				tmp := stack[len(stack)-2] - stack[len(stack)-1]
				stack[len(stack)-2] = tmp
				stack = stack[:len(stack)-1]
			case '*':
				tmp := stack[len(stack)-2] * stack[len(stack)-1]
				stack[len(stack)-2] = tmp
				stack = stack[:len(stack)-1]
			case '/':
				if stack[len(stack)-1] < 1e-10 && stack[len(stack)-1] > -1e-10 {
					panic("除数不能为0")
				}
				tmp := stack[len(stack)-2] / stack[len(stack)-1]
				stack[len(stack)-2] = tmp
				stack = stack[:len(stack)-1]
			}
		}
	}
	return stack[0]
}

//计算表达式.
func CalExpression(str string) float64 {
	var queue []*Node
	queue = ShuntingYard(str)
	res := CalNodes(queue)
	fmt.Printf("%v = %v \n", str, res)
	return res
}

func main() {
	str := " 3 + 3 * ( 6 - 1 * 2 ) / 4 - 1"
	result := CalExpression(str)
	fmt.Printf("The result: %f\n", result)
}

// 优化
// 可以构造一个符号优先级map，用于优先级比较
