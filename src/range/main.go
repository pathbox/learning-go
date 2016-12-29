for range的引入提升了Go的表达能力，但for range显然不是”免费的午餐“，在享用这个美味前，需要搞清楚for range的一些坑。

1、iteration variable重用

for range的idiomatic的使用方式是使用short variable declaration（:=）形式在for expression中声明iteration variable，但需要注意的是这些variable在每次循环体中都会被重用，而不是重新声明。

//details-in-go/5/iterationvariable.go
… …
    var m = [...]int{1, 2, 3, 4, 5}

    for i, v := range m {
        go func() {
            time.Sleep(time.Second * 3)
            fmt.Println(i, v)
        }()
    }

    time.Sleep(time.Second * 10)
… …

在我的Mac上，输出结果如下：

$go run iterationvariable.go
4 5
4 5
4 5
4 5
4 5

各个goroutine中输出的i,v值都是for range循环结束后的i, v最终值，而不是各个goroutine启动时的i, v值。一个可行的fix方法：

    for i, v := range m {
        go func(i, v int) {
            time.Sleep(time.Second * 3)
            fmt.Println(i, v)
        }(i, v)
    }

2、range expression副本参与iteration

range后面接受的表达式的类型包括：array, pointer to array, slice, string, map和channel(有读权限的)。我们以array为例来看一个简单的例子：

//details-in-go/5/arrayrangeexpression.go
func arrayRangeExpression() {
    var a = [5]int{1, 2, 3, 4, 5}
    var r [5]int

    fmt.Println("a = ", a)

    for i, v := range a {
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }
        r[i] = v
    }

    fmt.Println("r = ", r)
}

我们期待输出结果：

a =  [1 2 3 4 5]
r =  [1 12 13 4 5]
a =  [1 12 13 4 5]

但实际输出结果却是：

a =  [1 2 3 4 5]
r =  [1 2 3 4 5]
a =  [1 12 13 4 5]

我们原以为在第一次iteration，也就是i = 0时，我们对a的修改(a[1] = 12，a[2] = 13)会在第二次、第三次循环中被v取出，但结果却是v取出的依旧是a被修改前的值：2和3。这就是for range的一个不大不小的坑：range expression副本参与循环。也就是说在上面这个例子里，真正参与循环的是a的副本，而不是真正的a，伪代码如 下：

    for i, v := range a' {//a' is copy from a
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }
        r[i] = v
    }

Go中的数组在内部表示为连续的字节序列，虽然长度是Go数组类型的一部分，但长度并不包含的数组的内部表示中，而是由编译器在编译期计算出 来。这个例子中，对range表达式的拷贝，即对一个数组的拷贝，a'则是Go临时分配的连续字节序列，与a完全不是一块内存。因此无论a被 如何修改，其副本a'依旧保持原值，并且参与循环的是a'，因此v从a'中取出的仍旧是a的原值，而非修改后的值。

我们再来试试pointer to array：

func pointerToArrayRangeExpression() {
    var a = [5]int{1, 2, 3, 4, 5}
    var r [5]int

    fmt.Println("pointerToArrayRangeExpression result:")
    fmt.Println("a = ", a)

    for i, v := range &a {
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }

        r[i] = v
    }

    fmt.Println("r = ", r)
    fmt.Println("a = ", a)
    fmt.Println("")
}

这回的输出结果如下：

pointerToArrayRangeExpression result:
a =  [1 2 3 4 5]
r =  [1 12 13 4 5]
a =  [1 12 13 4 5]

我们看到这次r数组的值与最终a被修改后的值一致了。这个例子中我们使用了*[5]int作为range表达式，其副本依旧是一个指向原数组 a的指针，因此后续所有循环中均是&a指向的原数组亲自参与的，因此v能从&a指向的原数组中取出a修改后的值。

idiomatic go建议我们尽可能的用slice替换掉array的使用，这里用slice能否实现预期的目标呢？我们来试试：

func sliceRangeExpression() {
    var a = [5]int{1, 2, 3, 4, 5}
    var r [5]int

    fmt.Println("sliceRangeExpression result:")
    fmt.Println("a = ", a)

    for i, v := range a[:] {
        if i == 0 {
            a[1] = 12
            a[2] = 13
        }

        r[i] = v
    }

    fmt.Println("r = ", r)
    fmt.Println("a = ", a)
    fmt.Println("")
}

pointerToArrayRangeExpression result:
a =  [1 2 3 4 5]
r =  [1 12 13 4 5]
a =  [1 12 13 4 5]

显然用slice也能实现预期要求。我们可以分析一下slice是如何做到的。slice在go的内部表示为一个struct，由(*T, len, cap)组成，其中*T指向slice对应的underlying array的指针，len是slice当前长度，cap为slice的最大容量。当range进行expression复制时，它实际上复制的是一个 slice，也就是那个struct。副本struct中的*T依旧指向原slice对应的array，为此对slice的修改都反映到 underlying array a上去了，v从副本struct中*T指向的underlying array中获取数组元素，也就得到了被修改后的元素值。

slice与array还有一个不同点，就是其len在运行时可以被改变，而array的len是一个常量，不可改变。那么len变化的 slice对for range有何影响呢？我们继续看一个例子：

func sliceLenChangeRangeExpression() {
    var a = []int{1, 2, 3, 4, 5}
    var r = make([]int, 0)

    fmt.Println("sliceLenChangeRangeExpression result:")
    fmt.Println("a = ", a)

    for i, v := range a {
        if i == 0 {
            a = append(a, 6, 7)
        }

        r = append(r, v)
    }

    fmt.Println("r = ", r)
    fmt.Println("a = ", a)
}

输出结果：

a =  [1 2 3 4 5]
r =  [1 2 3 4 5]
a =  [1 2 3 4 5 6 7]

在这个例子中，原slice a在for range过程中被附加了两个元素6和7，其len由5增加到7，但这对于r却没有产生影响。这里的原因就在于a的副本a'的内部表示struct中的 len字段并没有改变，依旧是5，因此for range只会循环5次，也就只获取a对应的underlying数组的前5个元素。

range的副本行为会带来一些性能上的消耗，尤其是当range expression的类型为数组时，range需要复制整个数组；而当range expression类型为pointer to array或slice时，这个消耗将小得多，仅仅需要复制一个指针或一个slice的内部表示（一个struct）即可。我们可以通过 benchmark test来看一下三种情况的消耗情况对比：

对于元素个数为100的int数组或slice，测试结果如下：

//details-in-go/5/arraybenchmark
go test -bench=.
testing: warning: no tests to run
PASS
BenchmarkArrayRangeLoop-4             20000000           116 ns/op
BenchmarkPointerToArrayRangeLoop-4    20000000            64.5 ns/op
BenchmarkSliceRangeLoop-4             20000000            70.9 ns/op

可以看到range expression类型为slice或pointer to array的性能相近，消耗都近乎是数组类型的1/2。

3、其他range expression类型

对于range后面的其他表达式类型，比如string, map, channel，for range依旧会制作副本。

【string】
对string来说，由于string的内部表示为struct {*byte, len)，并且string本身是immutable的，因此其行为和消耗和slice expression类似。不过for range对于string来说，每次循环的单位是rune(code point的值)，而不是byte，index为迭代字符码点的第一个字节的position：

    var s = "中国人"

    for i, v := range s {
        fmt.Printf("%d %s 0x%x\n", i, string(v), v)
    }

输出结果：
0 中 0x4e2d
3 国 0x56fd
6 人 0x4eba

如果s中存在非法utf8字节序列，那么v将返回0xFFFD这个特殊值，并且在接下来一轮循环中，v将仅前进一个字节：

//byte sequence of s: 0xe4 0xb8 0xad 0xe5 0x9b 0xbd 0xe4 0xba 0xba
    var sl = []byte{0xe4, 0xb8, 0xad, 0xe5, 0x9b, 0xbd, 0xe4, 0xba, 0xba}
    for _, v := range sl {
        fmt.Printf("0x%x ", v)
    }
    fmt.Println("\n")

    sl[3] = 0xd0
    sl[4] = 0xd6
    sl[5] = 0xb9

    for i, v := range string(sl) {
        fmt.Printf("%d %x\n", i, v)
    }

输出结果：

0xe4 0xb8 0xad 0xe5 0x9b 0xbd 0xe4 0xba 0xba

0 4e2d
3 fffd
4 5b9
6 4eba

以上例子源码在details-in-go/5/stringrangeexpression.go中可以找到。

【map】

对于map来说，map内部表示为一个指针，指针副本也指向真实map，因此for range操作均操作的是源map。

for range不保证每次迭代的元素次序，对于下面代码：

 var m = map[string]int{
        "tony": 21,
        "tom":  22,
        "jim":  23,
    }

    for k, v := range m {
        fmt.Println(k, v)
    }

输出结果可能是：

tom 22
jim 23
tony 21

也可能是：

tony 21
tom 22
jim 23

或其他可能。

如果map中的某项在循环到达前被在循环体中删除了，那么它将不会被iteration variable获取到。
    counter := 0
    for k, v := range m {
        if counter == 0 {
            delete(m, "tony")
        }
        counter++
        fmt.Println(k, v)
    }
    fmt.Println("counter is ", counter)

反复运行多次，我们得到的两个结果：

tony 21
tom 22
jim 23
counter is  3

tom 22
jim 23
counter is  2

如果在循环体中新创建一个map元素项，那该项元素可能出现在后续循环中，也可能不出现：

    m["tony"] = 21
    counter = 0

    for k, v := range m {
        if counter == 0 {
            m["lucy"] = 24
        }
        counter++
        fmt.Println(k, v)
    }
    fmt.Println("counter is ", counter)

执行结果：

tony 21
tom 22
jim 23
lucy 24
counter is  4

or

tony 21
tom 22
jim 23
counter is  3

以上代码可以在details-in-go/5/maprangeexpression.go中可以找到。

【channel】

对于channel来说，channel内部表示为一个指针，channel的指针副本也指向真实channel。

for range最终以阻塞读的方式阻塞在channel expression上（即便是buffered channel，当channel中无数据时，for range也会阻塞在channel上），直到channel关闭：

//details-in-go/5/channelrangeexpression.go
func main() {
    var c = make(chan int)

    go func() {
        time.Sleep(time.Second * 3)
        c <- 1
        c <- 2
        c <- 3
        close(c)
    }()

    for v := range c {
        fmt.Println(v)
    }
}

运行结果：

1
2
3

如果channel变量为nil，则for range将永远阻塞。

六、select求值

golang引入的select为我们提供了一种在多个channel间实现“多路复用”的一种机制。select的运行机制这里不赘述，但select的case expression的求值顺序我们倒是要通过一个例子来了解一下：

// details-in-go/6/select.go

func takeARecvChannel() chan int {
    fmt.Println("invoke takeARecvChannel")
    c := make(chan int)

    go func() {
        time.Sleep(3 * time.Second)
        c <- 1
    }()

    return c
}

func getAStorageArr() *[5]int {
    fmt.Println("invoke getAStorageArr")
    var a [5]int
    return &a
}

func takeASendChannel() chan int {
    fmt.Println("invoke takeASendChannel")
    return make(chan int)
}

func getANumToChannel() int {
    fmt.Println("invoke getANumToChannel")
    return 2
}

func main() {
    select {
    //recv channels
    case (getAStorageArr())[0] = <-takeARecvChannel():
        fmt.Println("recv something from a recv channel")

        //send channels
    case takeASendChannel() <- getANumToChannel():
        fmt.Println("send something to a send channel")
    }
}

运行结果：

$go run select.go
invoke takeARecvChannel
invoke takeASendChannel
invoke getANumToChannel

invoke getAStorageArr
recv something from a recv channel

通过例子我们可以看出：
1) select执行开始时，首先所有case expression的表达式都会被求值一遍，按语法先后次序。

invoke takeARecvChannel
invoke takeASendChannel
invoke getANumToChannel

例外的是recv channel的位于赋值等号左边的表达式（这里是：(getAStorageArr())[0]）不会被求值。

2) 如果选择要执行的case是一个recv channel，那么它的赋值等号左边的表达式会被求值：如例子中当goroutine 3s后向recvchan写入一个int值后，select选择了recv channel执行，此时对=左侧的表达式 (getAStorageArr())[0] 开始求值，输出“invoke getAStorageArr”。