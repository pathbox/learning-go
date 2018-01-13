package main

const (
	max     = 500000
	bufsize = 100
)

func test(data chan int, done chan struct{}) int {
	count := 0

	go func() {
		for x := range data {
			count += x
		}

		close(done)
	}()

	for i := 0; i < max; i++ {
		data <- i
	}

	close(data)

	<-done
	return count
}

func main() {
	data := make(chan int, bufsize)
	done := make(chan struct{})

	println(test(data, done))
}

const block = 500

func testBlock(data chan [block]int, done chan struct{}) int {
	count := 0

	go func() {
		for a := range data {
			for _, x := range a {
				count += x
			}
		}

		close(done)
	}()

	for i := 0; i < max; i += block {
		//按块打包
		var b [block]int // 数组
		for n := 0; n < block; n++ {
			b[n] = i + n
			if i+n == max-1 {
				break
			}
		}

		data <- b
	}

	close(data)

	<-done
	return count
}

func testSliceBlock(data chan []int, done chan struct{}) int {
	count := 0

	go func() {
		for a := range data {
			for _, x := range a {
				count += x
			}
		}

		close(done)
	}()

	for i := 0; i < max; i += block {
		//按块打包
		b := make([]int, block) // slice
		for n := 0; n < block; n++ {
			b[n] = i + n
			if i+n == max-1 {
				break
			}
		}

		data <- b
	}

	close(data)

	<-done
	return count
}
