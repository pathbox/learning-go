package main

/*
static int add(int a, int b) {
	return a+b;
}
*/
import "C"
import "fmt"

func main() {
	v, err := C.add(1, 1)
	fmt.Println(v, err)
}
