package main

import (
	"fmt"
)

type Calculator struct {
	acc float64
}

type opfunc func(float64, float64) float64

func (c *Calculator) Do(op opfunc, v float64) float64 {
	c.acc = op(c.acc, v)
	return c.acc
}

func Add(a, b float64) float64 {
	return a + b
}

func Sub(a, b float64) float64 {
	return a - b
}

func Mul(a, b float64) float64 {
	return a * b
}

func main() {
	var c Calculator
	fmt.Println(c.Do(Add, 100))
	fmt.Println(c.Do(Sub, 10))
	fmt.Println(c.Do(Mul, 10))
}
