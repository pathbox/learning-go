package main

import (
	"net/http"
)

func main() {
	http.ListenAndServe(":12345", http.FileServer(http.Dir(".")))
}

// 大部分的服务器逻辑都需要使用者编写对应的 Handler，不过有些 Handler 使用频繁，因此 net/http 提供了它们的实现。比如负责文件 hosting 的 FileServer、负责 404 的NotFoundHandler 和 负责重定向的RedirectHandler。下面这个简单的例子，把当前目录所有文件 host 到服务端：
