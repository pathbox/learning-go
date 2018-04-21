package main

import (
	"encoding/json"
	"testing"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
)

var bb []byte = []byte("{\"conn_id\":\"12345678901234567891234567890111\",\"category\":\"agent_state\"}")

func BenchmarkJsoniter(b *testing.B) {
	m := make(map[string]string)

	for n := 0; n < b.N; n++ {
		jsoniter.Unmarshal(bb, &m)

	}
}

func BenchmarkJson(b *testing.B) {
	m := make(map[string]string)

	for n := 0; n < b.N; n++ {
		json.Unmarshal(bb, &m)

	}
}

func BenchmarkJsonparser(b *testing.B) {

	for n := 0; n < b.N; n++ {
		jsonparser.GetString(bb, "conn_id")
		jsonparser.GetString(bb, "conn_id")

	}
}
