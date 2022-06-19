package main

import (
	"fmt"
)


func Sum1[T int|float64](a,b T) T {
	return a + b
}

type Slice1 [T int|float64|string] []T
type Map1 [KEY int|string, VALUE string|float64] map[KEY]VALUE

type Struct1 [T string|int|float64] struct {
	Title string 
	Content T 
}


var MySlice1 Slice1[int] = []int{1,2,3}

type MyStruct2[S int | string, P map[S]string] struct {
	Name    string
	Content S
	Job     P
}


//切片泛型
type Slice11[T int | string] []T

//结构体泛型，它的第二个泛型参数的类型是第一个切片泛型。
type Struct11[P int | string, V Slice11[P]] struct {
	Name  P
	Title V
}

type MyNumber interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func ForeachInt[T MyNumber] (list []T) {
	for _, t := range list {
		fmt.Println(t)
	}
}

type myInt interface {
    int | int8 | int16 | int32 | int64
}

type myUint interface {
    uint | uint8 | uint16 | uint32
}

type myFloat interface {
    float32 | float64
}

type myNumber interface {
	myInt | myUint | myFloat | string
}

func ForeachI[T myNumber](list []T) {
	for _, t := range list {
		fmt.Println(t)
	}
}


// type myInt interface {
//     int | int8 | int16 | int32 | int64
// }

// type myInt2 interface {
//     int | int64
// }

// type myFloat interface {
//     float32 | float64
// }

// //每一个自定义约束类型单独一行
// type myNumber interface {
// 	myInt
// 	myInt2
// }
// 这样，myNumber的约束类型就是取的是myInt和myInt2的交接，即myNumber的约束范围是：int|int64。那如果是2个没有交集的约束呢？我们如果用编辑器，编辑器就会提示提示错误了，提示是个它是空的约束，传任何类型都不行。因为go里面的任何值类型都不是空集，都是有类型的. Cannot use int as the type myNumber2 Type does not implement constraint 'myNumber2' because constraint type set is empty


/*
//申明1个约束范围
type IntAll interface {
	int | int64 | int32
}

//定义1个泛型切片
type MySliceInt[T IntAll] []T

//正确:
var MyInt1 MySliceInt[int]

//自定义一个int型的类型
type YourInt int

//错误：实例化会报错
var MyInt2 MySliceInt[YourInt]

我们运行后，会发现，第二个会报错，因为MySliceInt允许的是int作为类型实参，而不是YourInt, 虽然YourInt类型底层类型是int ，但它依旧不是 int类型）。

这个时候~就排上用处了，我们可以这样写就可以了，表示底层的超集类型。

type IntAll interface {
	~int | ~int64 | ~int32  //  只要底层是 int int64 int32类型就行
}
*/


//自定义一个类型约束
type Number interface{
	int | int32 | int64 | float64 | float32 
}


//定义一个泛型结构体，表示堆栈
type Stack[V Number] struct {
	size  int
	value []V
}

//加上Push方法
func (s *Stack[V]) Push(v V) {
	s.value = append(s.value, v)
	s.size++
}

//加上Pop方法
func (s *Stack[V]) Pop() V {
	e := s.value[s.size-1]
	if s.size != 0 {
		s.value = s.value[:s.size-1]
		s.size--
	}
	return e
}

type error interface {
	Error() string
}


type DemoNumber interface {
	int | float64
}

type MyInterface[T int | string] interface {
	WriteOne(data T) T
	ReadOne() T
}

type Note struct {
}

func (n Note) WriteOne(one string) string {
	return "hello"
}

func (n Note) ReadOne() string {
	return "good"
}

type MyInterface2[T int | string] interface {
	// int|string 接口包含约束元素int和string，只能用作类型参数。

	WriteOne(data T) T
	ReadOne() T
}

func main() {
	MySlice2 := Slice1[int]{1,2,3}
	myMap1 := Map1[int, string] {
		1: "hello",
		2: "good",
	}
	fmt.Println(Sum1[int](1,2))
	fmt.Println(Sum1[float64](1.1,2.2))
	fmt.Println(MySlice2)
	fmt.Println(myMap1)

	var myStruct1 Struct1[float64]

	myStruct1.Title = "hello"
	myStruct1.Content = 3.14

	fmt.Println(myStruct1)

	myStruct2 := Struct1[string] {
		Title: "hello",
		Content: "good",
	}
	fmt.Println(myStruct2)
	// go无法识别这个匿名写法，不支持匿名泛型结构体

	var MyStruct1 = MyStruct2[int, map[int]string]{
		Name:    "small",
		Content: 1,
		Job:     map[int]string{1: "ss"},
	}

	fmt.Printf("%+v", MyStruct1)

	var MyStruct2 = MyStruct2[string, map[string]string]{
		Name:    "small",
		Content: "yang",
		Job:     map[string]string{"aa": "ss"},
	}
	
	fmt.Printf("%+v", MyStruct2)


	myStruct3 := Struct11[int, Slice11[int]] {
		Name: 123,
		Title: []int{1,2,3},
	}

	fmt.Println(myStruct3)


	s1 := &Stack[int]{}

	//入栈
	s1.Push(1)
	s1.Push(2)
	s1.Push(3)
	fmt.Println(s1.size, s1.value)  // 3 [1 2 3]

	//出栈
	fmt.Println(s1.Pop())  //3
	fmt.Println(s1.Pop())  //2
	fmt.Println(s1.Pop())  //1

	//实例化成一个float64型的结构体堆栈
	s2 := &Stack[float64]{}
	s2.Push(1.1)
	s2.Push(2.2)
	s2.Push(3.3)
	fmt.Println(s2.Pop())  //3.3
	fmt.Println(s2.Pop())  //2.2
	fmt.Println(s2.Pop())  //1.1

	var one MyInterface[string] = Note{}
	fmt.Println(one.WriteOne("hello"))
	fmt.Println(one.ReadOne())
}