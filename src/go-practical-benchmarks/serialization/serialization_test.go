package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/proto"
)

type Book struct {
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	Pages    int      `json:"num_pages"`
	Chapters []string `json:"chapters"`
}

/*
syntax = "proto2";
package main;

message BookProto {
  required string title = 1;
  required string author = 2;
  optional int64 pages = 3;
  repeated string chapters = 4;
}
*/
// protoc --go_out=. book.proto

/*
//go:generate msgp -tests=false
type BookDef struct {
	Title    string   `msg:"title"`
	Author   string   `msg:"author"`
	Pages    int      `msg:"num_pages"`
	Chapters []string `msg:"chapters"`
}
*/
// go generate

func generateObject() *Book {
	return &Book{
		Title:    "The Art of Computer Programming, Vol. 2",
		Author:   "Donald E. Knuth",
		Pages:    784,
		Chapters: []string{"Random numbers", "Arithmetic"},
	}
}

func generateMessagePackObject() *BookDef {
	obj := generateObject()
	return &BookDef{
		Title:    obj.Title,
		Author:   obj.Author,
		Pages:    obj.Pages,
		Chapters: obj.Chapters,
	}
}

func generateProtoBufObject() *BookProto {
	obj := generateObject()
	return &BookProto{
		Title:    proto.String(obj.Title),
		Author:   proto.String(obj.Author),
		Pages:    proto.Int64(int64(obj.Pages)),
		Chapters: obj.Chapters,
	}
}

func BenchmarkJSONMarshal(b *testing.B) {
	obj := generateObject()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	out, err := json.Marshal(generateObject())
	if err != nil {
		panic(err)
	}

	obj := &Book{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err = json.Unmarshal(out, obj)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoBufMarshal(b *testing.B) {
	obj := generateProtoBufObject()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := proto.Marshal(obj)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoBufUnmarshal(b *testing.B) {
	out, err := proto.Marshal(generateProtoBufObject())
	if err != nil {
		panic(err)
	}

	obj := &BookProto{}

	b.ResetTimer() // 重新计算benchmark时间，上面代码不再加入benchmark计算
	for n := 0; n < b.N; n++ {
		err = proto.Unmarshal(out, obj)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMessagePackMarshal(b *testing.B) {
	obj := generateMessagePackObject()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := obj.MarshalMsg(nil)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMessagePackUnmarshal(b *testing.B) {
	obj := generateMessagePackObject()
	msg, err := obj.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}

	obj = &BookDef{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err = obj.UnmarshalMsg(msg)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobMarshal(b *testing.B) {
	obj := generateObject()

	enc := gob.NewEncoder(ioutil.Discard)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := enc.Encode(obj)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGobUnmarshal(b *testing.B) {
	obj := generateObject()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		err = enc.Encode(obj)
		if err != nil {
			panic(err)
		}
	}

	dec := gob.NewDecoder(&buf)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err = dec.Decode(&Book{})
		if err != nil {
			panic(err)
		}
	}
}

/*
BenchmarkJSONMarshal-4            1,000,000  1,239 ns/op  640 B/op   3 allocs/op
BenchmarkJSONUnmarshal-4            500,000  3,249 ns/op  432 B/op   8 allocs/op
BenchmarkProtoBufMarshal-4        3,000,000    504 ns/op  552 B/op   5 allocs/op
BenchmarkProtoBufUnmarshal-4      2,000,000    692 ns/op  432 B/op  10 allocs/op
BenchmarkMessagePackMarshal-4    10,000,000    134 ns/op  160 B/op   1 allocs/op
BenchmarkMessagePackUnmarshal-4   5,000,000    252 ns/op  112 B/op   4 allocs/op
BenchmarkGobMarshal-4             2,000,000    737 ns/op   32 B/op   1 allocs/op
BenchmarkGobUnmarshal-4           1,000,000  1,005 ns/op  272 B/op   8 allocs/op
*/
