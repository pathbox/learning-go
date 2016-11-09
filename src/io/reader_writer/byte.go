package main

import (
	"bytes"
	"fmt"
	"os"
)

func ByteRWerExample() {
  FOREND:
    for{
      fmt.Println("请输入要通过WriteByte写入的一个ASCII字符（b：返回上级；q：退出）：")
      var ch byte
      fmt.Scanf("%c\n", &ch)
      switch ch {
      case 'b':
        fmt.Println("返回上一级菜单！")
        break FOREND
      case 'q':
        fmt.Println("程序退出！")
        os.Exit(0)
      default:
        buffer := new(bytes.Buffer)
        err := buffer.WriteByte(ch)
      }
    }
}
