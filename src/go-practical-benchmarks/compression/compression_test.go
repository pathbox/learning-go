package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	data, err := ioutil.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	zw := gzip.NewWriter(ioutil.Discard)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err = zw.Write(data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	data, err := ioutil.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err = zw.Write(data)
	if err != nil {
		panic(err)
	}

	err = zw.Close()
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(buf.Bytes())
	zr, _ := gzip.NewReader(r)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Reset(buf.Bytes())
		zr.Reset(r)
		_, err := ioutil.ReadAll(zr)
		if err != nil {
			panic(err)
		}
	}
}
