package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsChinese(t *testing.T) {
	Convey("Init a Chinese string", t, func() {
		str := "测试"
		Convey("Call isChinese func", func() {
			result := isChinese(str)

			Convery("result should be: true", func() {
				So(result, ShouldEqual, true)
			})
		})
	})
}

func TestParseArgs(t *testing.T) {
	Convey("Init args arrays: withVoice & withOutVoice", t, func() {
		withVoice := []string{"aa", "bb", "-v"}
		withOutVoice := []string{"aa", "bb", "cc"}

		Convey("Call parse func", func() {
			_, ret := parseArgs(withVoice)
			_, ret2 := parseArgs(withOutVoice)

			Convey("result should be: true & false", func() {
				So(ret, ShouldEqual, true)
				So(ret2, ShouldEqual, false)
			})
		})
	})
}
