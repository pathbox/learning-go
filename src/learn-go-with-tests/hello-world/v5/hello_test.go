package main

import "testing"

func TestHello(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want string) { // 定义一个闭包，方便使用
		t.Helper()
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris")  // 调用要测试的方法，得到方法的返回值
		want := "Hello, Chris" // 设置想要的正确值

		assertCorrectMessage(t, got, want) // 将两个值进行比较
	})

	t.Run("empty string defaults to 'world'", func(t *testing.T) {
		got := Hello("")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})
}
