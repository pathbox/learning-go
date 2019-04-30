package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var types = []string{
	"int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"byte", "rune", "uintptr", "bool", "string",
	"float32", "float64", "complex64", "complex128",
	"[]byte", "[]string", "map[string]int",
	"chan int", "func(int) int",
}

var values = []string{
	"0", "0", "0", "0", "0",
	"0", "0", "0", "0", "0",
	"0", "0", "0", "false", "\"\"",
	"0", "0", "0+0i", "0+0i",
	"[]byte{}", "[]string{}", "map[string]int{}",
	"nil", "func(int) int {return 0}",
}

const template1 = `package main

import (
	"fmt"
	"reflect"
)

var v = struct {
		a %v
		b %v
		c %v
		d %v
		e %v
}{%v, %v, %v, %v, %v}
`
const template2 = `
func init() {
	fmt.Printf("%#T\n", v)
    t := reflect.TypeOf(v)
    fmt.Printf("结构体大小：%v\n", t.Size())
    for i := 0; i < t.NumField(); i++ {
        showAlign(t, i)
    }
}

func showAlign(v reflect.Type, i int) {
    sf := v.Field(i)
    fmt.Printf("字段 %10v，大小：%2v，对齐：%2v，字段对齐：%2v，偏移：%2v\n",
        sf.Type.Kind(),
        sf.Type.Size(),
        sf.Type.Align(),
        sf.Type.FieldAlign(),
        sf.Offset,
    )
}`

func main() {
	f, err := os.OpenFile("testAlign.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	t := [5]string{}
	v := [5]string{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		n := rand.Intn(len(types))
		t[i] = types[n]
		v[i] = values[n]
	}

	fmt.Fprintf(f, template1,
		t[0], t[1], t[2], t[3], t[4],
		v[0], v[1], v[2], v[3], v[4],
	)
	fmt.Fprint(f, template2)
}
