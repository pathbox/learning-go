package ydict

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	client := NewOnlineClient("go-ydict", "252639882")
	if client.BaseURL != "http://fanyi.youdao.com/" {
		t.Error("res.BaseURL should equal", "http://dict.youdao.com/")
	}
	fmt.Printf("%+v\n", client)
	res, err := client.Query("hello")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", res)

}
