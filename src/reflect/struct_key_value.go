package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type User struct {
	ID        int64     `json:"id"`         // 自增主键
	Age       int64     `json:"age"`        // 年龄
	FirstName string    `json:"first_name"` // 姓
	LastName  string    `json:"last_name"`  // 名
	Email     string    `json:"email"`      // 邮箱地址
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

func StructKeyValue(v reflect.Value) ([]string, []string) {
	var keys, values []string
	t := v.Type()

	for n := 0; n < t.NumField(); n++ {
		tf := t.Field(n)
		vf := v.Field(n)

		fmt.Println("tf===:", tf)
		fmt.Println("vf===:", vf)
		// 忽略非导出字段
		if tf.Anonymous {
			continue
		}
		// 忽略无效字段 零值字段
		if !vf.IsValid() || reflect.DeepEqual(vf.Interface(), reflect.Zero(vf.Type()).Interface()) {
			continue
		}
		for vf.Type().Kind() == reflect.Ptr {
			vf = vf.Elem()
		}
		//有时候根据需求会组合struct，这里处理下，支持获取嵌套的struct tag和value
		//如果字段值是time类型之外的struct，递归获取keys和values
		if vf.Kind() == reflect.Struct && tf.Type.Name() != "Time" {
			cKeys, cValues := StructKeyValue(vf)
			keys = append(keys, cKeys...)
			values = append(values, cValues...)
			continue
		}
		fmt.Println("tag:", tf.Tag.Get("json"))
		key := strings.Split(tf.Tag.Get("json"), ",")[0]
		if key == "" {
			continue
		}
		value := format(vf)
		if value != "" {
			keys = append(keys, key)
			values = append(values, value)
		}
	}
	return keys, values
}

func format(v reflect.Value) string {
	//断言出time类型直接转unix时间戳
	if t, ok := v.Interface().(time.Time); ok {
		return fmt.Sprintf("FROM_UNIXTIME(%d)", t.Unix())
	}
	switch v.Kind() {
	case reflect.String:
		return fmt.Sprintf(`'%s'`, v.Interface())
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return fmt.Sprintf(`%d`, v.Interface())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(`%f`, v.Interface())
	//如果是切片类型，遍历元素，递归格式化成"(, , , )"形式
	case reflect.Slice:
		var values []string
		for i := 0; i < v.Len(); i++ {
			values = append(values, format(v.Index(i)))
		}
		return fmt.Sprintf(`(%s)`, strings.Join(values, ","))
	//接口类型剥一层递归
	case reflect.Interface:
		return format(v.Elem())
	}
	return ""
}

func MapKeyValue(v reflect.Value) ([]string, []string) {
	var keys, values []string

	mapKeys := v.MapKeys()
	for _, key := range mapKeys {
		fmt.Println("map value:", v.MapIndex(key))
		value := format(v.MapIndex(key))
		if value != "" {
			values = append(values, value)
			keys = append(keys, key.Interface().(string))
		}
	}
	return keys, values
}

func main() {
	user := &User{
		ID:        1,
		Age:       28,
		FirstName: "Cary",
		LastName:  "Cary",
		Email:     "Cary@email.com",
	}

	uValue := reflect.ValueOf(user)
	for uValue.Kind() == reflect.Ptr {
		uValue = uValue.Elem()
	}

	ks, vs := StructKeyValue(uValue)

	fmt.Println("Keys:", ks)
	fmt.Println("Values:", vs)

	uMap := map[string]string{
		"Name": "Cary",
		"City": "New York",
	}
	umValue := reflect.ValueOf(uMap)
	mks, mvs := MapKeyValue(umValue)
	fmt.Println("mKeys:", mks)
	fmt.Println("mValues:", mvs)
}
