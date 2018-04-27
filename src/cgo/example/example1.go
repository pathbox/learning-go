package main

// #include <stdio.h> // 启用cgo
import "C"

func main() {
	C.puts(C.CString("你好， Hello World"))
}
