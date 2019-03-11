package main

import "fmt"

func quickSort(nums []int, ch chan int, level int, threads int) {
	level = level * 2
	if len(nums) == 1 { // 数组只剩一个整数，将该整数写到ch
		ch <- nums[0]
		close(ch)
		return
	}

	if len(nums) == 0 {
		close(ch)
		return
	}

	less := make([]int, 0)
	greater := make([]int, 0)
	left := nums[0] // 每次的基准值
	nums = nums[1:] // 基准值之后的数据

	for _, numData := range nums {
		switch {
		case numData <= left:
			less = append(less, numData)
		case numData > left:
			greater = append(greater, numData)
		}
	}

	leftCh := make(chan int, len(less))
	rightCh := make(chan int, len(greater))

	if level <= threads {
		go quickSort(less, leftCh, level, threads) //分任务  创建 goroutine 执行
		go quickSort(greater, rightCh, level, threads)
	} else {
		quickSort(less, leftCh, level, threads) // 没有创建goroutine 执行
		quickSort(greater, rightCh, level, threads)
	}

	// 合并数据

	for i := range leftCh {
		ch <- i
	}
	// fmt.Println("base", left)
	ch <- left // 将基准值传给ch

	for i := range rightCh {
		ch <- i
	}
	close(ch)
	return
}

func main() {
	x := []int{3, 1, 4, 1, 5, 11, 88, 23, 9, 2, 6}
	ch := make(chan int)
	go quickSort(x, ch, 0, 0) // 0 0 表示不限制线程个数
	for v := range ch {
		fmt.Println(v)
	}
}
