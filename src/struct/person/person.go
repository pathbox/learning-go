package person

// 大写表示别的包导入person时 可以是用这个struct， 是public的
type Person struct {
	Name string
	Age  int
}

// 小写表示是 private， 只能在当前的代码逻辑中使用
type person struct {
}
