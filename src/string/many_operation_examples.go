// 字符串长度
// 由于string类型可以看成是一种特殊的slice类型，因此获取长度可以用内置的函数len；同时支持 切片 操作，因此，子串获取很容易。

len("Hello World!")

import "strings"
// 是否存在某个字符或子串

// 子串substr在s中，返回true
func Contains(s, substr string) bool
// chars中任何一个Unicode代码点在s中，返回true
func ContainsAny(s, chars string) bool
// Unicode代码点r在s中，返回true
func ContainsRune(s string, r rune) bool

fmt.Println(strings.ContainsAny("team", "i"))
fmt.Println(strings.ContainsAny("failure", "u & i"))
fmt.Println(strings.ContainsAny("in failure", "s g"))
fmt.Println(strings.ContainsAny("foo", ""))
fmt.Println(strings.ContainsAny("", ""))

// 子串出现次数(字符串匹配)

// 在Go中，查找子串出现次数即字符串模式匹配，实现的是Rabin-Karp算法
func Count(s, sep string) int
// 在 Count 的实现中，处理了几种特殊情况，属于字符匹配预处理的一部分。这里要特别说明一下的是当 sep 为空时，Count 的返回值是：utf8.RuneCountInString(s) + 1

fmt.Println(strings.Count("five", "")) // before & after each rune

// 另外，Count 是计算子串在字符串中出现的无重叠的次数，比如：

fmt.Println(strings.Count("fivevev", "vev"))  // 输出: 1

// 字符串分割为[]string

func Fields(s string) []string
func FieldsFunc(s string, f func(rune) bool) []string

fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   ")) // Fields are: ["foo" "bar" "baz"]

// Split 和 SplitAfter、 SplitN 和 SplitAfterN
// 之所以将这四个函数放在一起讲，是因为它们都是通过一个同一个内部函数来实现的

func Split(s, sep string) []string { return genSplit(s, sep, 0, -1) }
func SplitAfter(s, sep string) []string { return genSplit(s, sep, len(sep), -1) }
func SplitN(s, sep string, n int) []string { return genSplit(s, sep, 0, n) }
func SplitAfterN(s, sep string, n int) []string { return genSplit(s, sep, len(sep), n) }

// 那么，Split 和 SplitAfter 有啥区别呢？通过这两句代码的结果就知道它们的区别了：
fmt.Printf("%q\n", strings.Split("foo,bar,baz", ","))   // ["foo" "bar" "baz"]
fmt.Printf("%q\n", strings.SplitAfter("foo,bar,baz", ",")) // ["foo," "bar," "baz"]

// 带 N 的方法可以通过最后一个参数 n 控制返回的结果中的 slice 中的元素个数，当 n < 0 时，返回所有的子字符串；当 n == 0 时，返回的结果是 nil；当 n > 0 时，表示返回的 slice 中最多只有 n 个元素，其中，最后一个元素不会分割，比如：

fmt.Printf("%q\n", strings.SplitN("foo,bar,baz", ",", 2)) // ["foo" "bar,baz"]

fmt.Printf("%q\n", strings.Split("a,b,c", ",")) // ["a" "b" "c"]
fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a ")) // ["" "man " "plan " "canal panama"]
fmt.Printf("%q\n", strings.Split(" xyz ", "")) // [" " "x" "y" "z" " "]
fmt.Printf("%q\n", strings.Split("", "Bernardo O'Higgins")) // [""]

// 字符串是否有某个前缀或后缀

func HasPrefix(s, prefix string) bool {
  return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func HasSuffix(s, suffix string) bool {
  return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// 字符或子串在字符串中出现的位置
// 在 s 中查找 sep 的第一次出现，返回第一次出现的索引
func Index(s, sep string) int
// chars中任何一个Unicode代码点在s中首次出现的位置
func IndexAny(s, chars string) int
// 查找字符 c 在 s 中第一次出现的位置，其中 c 满足 f(c) 返回 true
func IndexFunc(s string, f func(rune) bool) int
// Unicode 代码点 r 在 s 中第一次出现的位置
func IndexRune(s string, r rune) int

// 有三个对应的查找最后一次出现的位置
func LastIndex(s, sep string) int
func LastIndexAny(s, chars string) int
func LastIndexFunc(s string, f func(rune) bool) int

fmt.Printf("%d\n", strings.IndexFunc("studygolang", func(c rune) bool {
    if c > 'u' {
        return true
    }
    return false
}))

// 字符串 JOIN 操作
// 将字符串数组（或slice）连接起来可以通过 Join 实现

func Join(a []string, sep string) string

// 假如没有这个库函数，我们自己实现一个，我们会这么实现：
func Join(str []string, sep string) string {
  if len(str) == 0 {
    return ""
  }
  if len(str) == 1{
    return str[0]
  }
  buffer := bytes.NewBufferString(str[0])
  for _, s := range str[1:] {
    buffer.WriteString(sep)
    buffer.WriteString(s)
  }
  return buffer.String()
}

// 这里，我们使用了 bytes 包的 Buffer 类型，避免大量的字符串连接操作（因为 Go 中字符串是不可变的）

func Join(a []string, sep string) string{
  if len(a) == 0 {
    return ""
  }
  if len(a) == 1 {
    return a[0]
  }

  n := len(sep) * (len(a) -1)
  for i := 0; i < len(a); i++ {
    n += len(a[i])
  }

  b := make([]byte, n)
  bp := copy(b, a[0])
  for _, s := range a[1:] {
    bp += copy(b[bp:], sep)
    bp += copy(b[bp:], s)
  }
  return string(b)
}
// 标准库的实现没有用 bytes 包，当然也不会简单的通过 + 号连接字符串。Go 中是不允许循环依赖的，标准库中很多时候会出现代码拷贝，而不是引入某个包。这里 Join 的实现方式挺好，我个人猜测，不直接使用 bytes 包，也是不想依赖 bytes 包（其实 bytes 中的实现也是 copy 方式）。

// 简单使用示例：

fmt.Println(Join([]string{"name=xxx", "age=xx"}, "&"))
// 输出 name=xxx&age=xxx

// 字符串重复几次
func Repeat(s string, count int) string
fmt.Println("ba" + strings.Repeat("na", 2))

// 字符串子串替换
// 用 new 替换 s 中的 old，一共替换 n 个。
// 如果 n < 0，则不限制替换次数，即全部替换
func Replace(s, old, new string, n int) string

fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))  // oinky oinky oink

fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1)) // moo moo moo

// Replacer 类型

r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
fmt.Println(r.Replace("This is <b>HTML</b>!"))

 // Reader 类型
// 这是实现了 io 包中的接口。它实现了 io.Reader（Read 方法），io.ReaderAt（ReadAt 方法），io.Seeker（Seek 方法），io.WriterTo（WriteTo 方法），io.ByteReader（ReadByte 方法），io.ByteScanner（ReadByte 和 UnreadByte 方法），io.RuneReader（ReadRune 方法） 和 io.RuneScanner（ReadRune 和 UnreadRune 方法）。

type Reader struct {
    s        string // Reader 读取的数据来源
    i        int // current reading index（当前读的索引位置）
    prevRune int // index of previous rune; or < 0（前一个读取的 rune 索引位置）
}

func NewReader(s string) *Reader