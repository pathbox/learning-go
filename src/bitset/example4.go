package main

import (
	"fmt"

	"github.com/bits-and-blooms/bitset"
)

/*
假设现在数据库里有一个字段存储用户状态，设计是这样的：0 00 00

第1、2位表示会员等级 00表示普通会员，01表示vip1，10表示vip2，11表示svip

第3、4位表示头像状态 00表示未上传，01表示01审核中，10审核失败，11审核通过

第5位表示账号状状态 0表示正常，1表示封禁
*/

func main() {
	b := bitset.New(5)
	// 设置vip1，第1位0，第2位1
	b.SetTo(1, false).SetTo(2, true)

	// 设置头像审核失败，第3位1，第4位0
	b.SetTo(3, true).SetTo(4, false)

	// 状态初始化
	b.ClearAll()

	// 查看账号状态，第5位，true代表1 false代表0
	b.Test(5)

	// 是不是svip，第1、2位是11
	r := b.Test(1) && b.Test(2)
	fmt.Println(r)

}
