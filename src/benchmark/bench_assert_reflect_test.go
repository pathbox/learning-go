package benchmark

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Document Document
type Document struct {
	Statement []Statement `json:"Statement"`
}

// Statement Statement
type Statement struct {
	Action interface{} `json:"Action"`
}

func Benchmark_Assert(b *testing.B) {
	val := getInterface()
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	// b.StopTimer()
	// // 开始计时器
	// b.StartTimer()
	for i := 0; i < b.N; i++ {
		valueToStringSlice(val)
	}
}

func Benchmark_Reflect(b *testing.B) {
	val := getInterface()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReflectToSlice(val)
	}
}

func Benchmark_Parallel_Assert(b *testing.B) {
	val := getInterface()
	// 重置计时器
	b.ResetTimer()
	// 停止计时器
	// b.StopTimer()
	// // 开始计时器
	// b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			valueToStringSlice(val)
		}
	})
}

func Benchmark_Parallel_Reflect(b *testing.B) {
	val := getInterface()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ReflectToSlice(val)
		}
	})
}

func valueToStringSlice(value interface{}) []string {
	strSlice := make([]string, 0)
	switch v := value.(type) {
	case string:
		strSlice = append(strSlice, v)
	case []interface{}:
		for _, val := range v {
			if val != nil {
				strSlice = append(strSlice, val.(string))
			}
		}
	}
	return strSlice
}

func ReflectToSlice(value interface{}) []string {
	valueSlice := make([]string, 0)
	// 反射处理
	typeOfA := reflect.TypeOf(value)
	valueOfA := reflect.ValueOf(value)
	switch typeOfA.Kind() {
	case reflect.Slice:
		for i := 0; i < valueOfA.Len(); i++ {
			val := valueOfA.Index(i)
			valueSlice = append(valueSlice, val.Interface().(string))
		}
	case reflect.String:
		valueSlice = append(valueSlice, valueOfA.String())
	}
	return valueSlice
}

func unmarshalDoc(docStr string) (*Document, error) {
	docByte := []byte(docStr)
	doc := &Document{}
	err := json.Unmarshal(docByte, &doc)
	return doc, err
}

func getInterface() interface{} {
	jsonStr := `{
	    "Statement": [
	        {
	            "Effect": "Allow",
	            "Action": [
	                      "oss:ListBuckets",
	                      "oss:GetBucketStat",
	                      "oss:GetBucketInfo",
	                      "oss:GetBucketTagging",
	                      "oss:GetBucketAcl"
	                      ],
	            "Resource": [
	                "acs:oss:*:*:*"
	            ]
	        },
	        {
	            "Effect": "Allow",
	            "Action": [
	                "oss:GetObject",
	                "oss:GetObjectAcl"
	            ],
	            "Resource": [
	                "acs:oss:*:*:myphotos/hangzhou/2015/*"
	            ]
	        },
	        {
	            "Effect": "Allow",
	            "Action": "oss:ListObjects",
	            "Resource": [
	                "acs:oss:*:*:myphotos"
	            ],
	            "Condition": {
	                "StringLike": {
	                    "oss:Delimiter": "/",
	                    "oss:Prefix": [
	                        "",
	                        "hangzhou/",
	                        "hangzhou/2015/*"
	                    ]
	                }
	            }
	        }
	    ]
	}`
	doc, err := unmarshalDoc(jsonStr)
	if err != nil {
		panic(err)
	}
	st := doc.Statement[0]
	return st.Action
}

//
/* go test -v -bench=. -benchmem bench_assert_reflect_test.go
goos: darwin
goarch: amd64
Benchmark_Assert
Benchmark_Assert-4    	 3210492	       368 ns/op
Benchmark_Reflect
Benchmark_Reflect-4   	 2559997	       491 ns/op
PASS
ok  	command-line-arguments	3.307s

interface{} 值反射的方式比 断言的方式性能差
*/
