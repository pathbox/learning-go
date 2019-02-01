package main

import "fmt"

func main() {
	slice := []int{0, 1, 2, 3}
	myMap := make(map[int]*int)

	for index, value := range slice {
		fmt.Printf("value: %p\n", &value)
		// myMap[index] = &value //这时就会打出坑爹的结果
		num := value
		fmt.Printf("num: %p\n", &num)
		myMap[index] = &num
	}
	fmt.Println("=====new map=====")
	prtMap(myMap)
}

func prtMap(myMap map[int]*int) {
	for key, value := range myMap {
		fmt.Printf("map[%v]=%v\n", key, *value)
	}
}

/*
=====new map=====
map[3]=3
map[0]=3
map[1]=3
map[2]=3

由输出可以知道，不是我们预期的输出，正确输出应该如下：
=====new map=====
map[0]=0
map[1]=1
map[2]=2
map[3]=3

for range创建了每个元素的副本，而不是直接返回每个元素的引用，如果使用该值变量的地址作为指向每个元素的指针，就会导致错误，在迭代时，返回的变量是一个迭代过程中根据切片依次赋值的新变量，并且这个新变量的地址总是相同的



package main

import "fmt"

func main() {
    slice := []int{0, 1, 2, 3}
    myMap := make(map[int]*int)

    for index, value := range slice {
        fmt.Printf("value: %p\n", &value)
				num := value
				fmt.Printf("num: %p\n", &num)
				myMap[index] = &num
    }
    fmt.Println("=====new map=====")
    prtMap(myMap)
}

func prtMap(myMap map[int]*int) {
    for key, value := range myMap {
        fmt.Printf("map[%v]=%v\n", key, *value)
    }
}

这里是num每次循环新建的一个变量，所以地址不同
将地址赋值给myMap。
出了for循环后，myMap并没有回收，所以每次循环中创建的num也不会回收。

*/
