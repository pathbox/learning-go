package main

import "fmt"

var array = []int{99, 5, 36, 7, 22, 17, 46, 12, 2, 19, 25, 28, 1, 92}
var n = len(array) - 1

func main() {

	heapSort()

	fmt.Println("After build: ", array)
}

// 传入一个需要向下调整的节点编号i， i其实都为顶点节点
// 该方法只是从3个节点中进行调整， 父节点-右儿子-左儿子
func siftdown(i int) {
	var t, flag int // flag用来标记是否需要继续向下调整
	// 当i节点有儿子，并且有需要继续调整的时候循环就执行
	for i*2+1 <= n && flag == 0 {
		// 首先判断它和左儿子的关系，并用t记录值较小的节点编号
		if array[i] > array[i*2+1] {
			t = i*2 + 1
		} else {
			t = i
		}

		//如果它有右儿子，再对右儿子进行讨论
		if i*2+2 <= n {
			// 如果右儿子的值更小，更新较小的节点编号
			if array[t] > array[i*2+2] {
				t = i*2 + 2
			}
		}
		// 如果发现最小的节点编号不是自己，说明子节点中有比父节点更小的.将父节点和儿子节点交换
		if t != i {
			swap(t, i)
			i = t // 更新i为刚才与它交换的儿子节点编号，便于接下来继续向下调整
		} else {
			flag = 1 // 说明当前的父节点已经比两个子节点要小了，不需要进行调整
		}
	}
}

func swap(x, y int) {
	array[x], array[y] = array[y], array[x]
}

// 每次输出 最祖先顶点
func heapSort() {
	for n >= 0 {
		swap(0, n)
		n--
		siftdown(0)
	}
}
