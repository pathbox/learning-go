package main

import "fmt"

type Tester interface {
	Test(int)
}

type Data struct {
	x int
}

func (d *Data) Test(x int) {
	d.x += x
}

func call(d *Data) {
	d.Test(100)
}

func ifaceCall(t Tester) {
	t.Test(100)
}

func main() {
	d := &Data{x: 100}

	call(d)
	ifaceCall(d)

	fmt.Println(d)
}
