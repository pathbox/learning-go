https://zhuanlan.zhihu.com/p/25474088

Golang 反射慢的原因

```go
type_ := reflect.TypeOf(obj)
field, _ := type_.FieldByName("hello")

// 这里取出来的 field 对象是 reflect.StructField 类型，但是它没有办法用来取得对应对象上的值。如果要取值，得用另外一套对object，而不是type的反射

type_ := reflect.ValueOf(obj)
fieldValue := type_.FieldByName("hello")

// 每次反射都需要malloc这个reflect.Value结构体

```

<!-- Jsoniter json解析器 使用的原因是 用 reflect.Type 得出来的信息来直接做反射，而不依赖于 reflect.ValueOf -->