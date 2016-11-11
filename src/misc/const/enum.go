package main

// 关键字 iota 定义常量组中从 0 开始按行计数的自增枚举值

const (
	Sunday    = iota // 0
	Monday           // 1, 通常省略后续行行表达式
	Tuesday          // 2
	Wednesday        // 3
	Thursday         // 4
	Friday           // 5
	Saturday         // 6
)

const (
	_        = iota             // iota = 0
	KB int64 = 1 << (10 * iota) // iota = 1
	MB                          // 与 KB 表达式相同，但 iota = 2
	GB
	TB
)

// 在同一常量组中，可以提供多个 iota，它们各自增长。
const (
	A, B = iota, iota << 10 // 0, 0 << 10
	C, D                    // 1, 1 << 10
)

// 如果 iota 自增被打断，须显式恢复
const (
	A = iota // 0
	B        // 1
	C = "c"  // c
	D        // c, 与上⼀行相同
	E = iota // 4，显式恢复。注意计数包含了 C、D 两⾏行。
	F        // 5
)

func main() {

}
