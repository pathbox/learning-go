package main

func fib(a int) int {
	if a == 1 {
		return 0
	}

	if a == 2 {
		return 1
	}

	return fib(a-1) + fib(a-2)
}

func main() {
	println(fib(15))
}
