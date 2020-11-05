package main

import (
	"fmt"
	"path/filepath"

	"github.com/gobwas/glob"
)

func main() {
	var b bool
	var err error
	t1 := "ucs:oss:ListBuckets"
	p1 := "ucs:oss:*"

	b, err = filepath.Match(p1, t1)
	fmt.Println(b, err)

	t2 := "ucs:oss:ListBuckets"
	p2 := "ucs:*:List*"

	b, err = filepath.Match(p2, t2)
	fmt.Println(b, err)

	t3 := "ucs:oss:ListBuckets"
	p3 := "ucs:*:ListBuckets"

	b, err = filepath.Match(p3, t3)
	fmt.Println(b, err)

	t4 := "ucs:oss:ListBuckets"
	p4 := "ucs:oss:*Buckets"

	b, err = filepath.Match(p4, t4)
	fmt.Println(b, err)

	t5 := "ucs:oss:ListBuckets"
	p5 := "*"

	b, err = filepath.Match(p5, t5)
	fmt.Println(b, err)

	t6 := "urn:oss:b:200000456:myphotos/hangzhou/2015/aaa"
	p6 := "urn:oss:*:200000456:myphotos/hangzhou/2015/*"

	b, err = filepath.Match(p6, t6)
	fmt.Println(b, err)

	t7 := "urn:oss:22:200000456:myphotos"
	p7 := "urn:oss:*:200000456:myphotos"

	g, _ := glob.Compile(p7)
	fmt.Println("ggg:", g.Match(t7))
}
