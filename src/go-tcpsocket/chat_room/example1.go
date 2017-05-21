package main

import (
	"fmt"
	"net"
	"os"
)

func checkError(err error, info string) (res bool) {

	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}

////////////////////////////////////////////////////////
//
//服务器端接收数据线程
//参数：
//      数据连接 conn
//      通讯通道 messages
//
////////////////////////////////////////////////////////

// 每次有新的client连接产生时，执行一次
func Handler(conn net.Conn, message chan string) {
	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())
	fmt.Println("Handler doing")

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		defer conn.Close()
		if checkError(err, "Connection") == false {
			conn.Close()
			break
		}
		if length > 0 {
			buf[length] = 0
		}
		fmt.Println("Rec[", conn.RemoteAddr().String(), "] Say :", string(buf[0:length]))
		reciveStr := string(buf[0:length])
		message <- reciveStr
	}
}

////////////////////////////////////////////////////////
//
//服务器发送数据的线程
//
//参数
//      连接字典 conns
//      数据通道 messages
//
////////////////////////////////////////////////////////

func echoHandler(conns *map[string]net.Conn, message chan string) {
	for { // 死循环阻塞监听 <-message 的到来
		msg := <-message // 有消息通过channel 传过来了，则继续执行，然后又回到for的开始并阻塞在这里，等待新的message
		fmt.Println("msg: ", msg)
		// 每次有新的client消息到来时，执行一次 给所有client echo 消息
		for key, value := range *conns { // 这里的conns为什么会有值呢，在StartServer 中，第95行给conns赋值，创建新的goroutine处理Handler，在Handler中，是处理client发送的数据，读取到之后，传给 message channel。第58行会一直阻塞等待第42行的执行，当执行第42行的时候，conns已经在之前执行了，必定有值，这里conns传的是map的指针
			fmt.Println("connection is connected from ...", key)
			fmt.Println("echoHandler doing")
			_, err := value.Write([]byte(msg))
			defer value.Close()
			if err != nil {
				fmt.Println(err.Error())
				delete(*conns, key)
			}
		}
	}
}

////////////////////////////////////////////////////////
//
//启动服务器
//参数
//  端口 port
//
////////////////////////////////////////////////////////
func StartServer(port string) {
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	conns := make(map[string]net.Conn)
	messages := make(chan string, 10)
	//启动服务器广播线程
	go echoHandler(&conns, messages)

	for {
		fmt.Println("Listening ...") // 虽然这个在死循环for中，但其实只会在初始的时候执行一次，然后每次产生新的client连接后执行一次，因为l.Accept()
		conn, err := l.Accept()
		checkError(err, "Accept")
		fmt.Println("Accepting .......")
		conns[conn.RemoteAddr().String()] = conn
		go Handler(conn, messages) // 为什么 在死循环for中起goroutine，不会造成大量的goroutine产生吗？ 不会，当l.Accept()接收到一个conn的时候，才会继续往下执行，要不就会一直阻塞在l.Accept()这一步。 每个conn表示一个client，每个client对应一个Handler goroutine进行处理
	}
}

////////////////////////////////////////////////////////
//
//客户端发送线程
//参数
//      发送连接 conn
//
////////////////////////////////////////////////////////

func chatSend(conn net.Conn) {
	var input string
	username := conn.LocalAddr().String()

	for {
		fmt.Scanln(&input) // 从终端得到键盘输入的字符串
		if input == "/quit" {
			fmt.Println("ByeBye..")
			conn.Close()
			os.Exit(0)
		}

		lens, err := conn.Write([]byte(username + " Say --:" + input)) // 将输入的字符串 通过连接conn 发送给服务端
		fmt.Println(lens)
		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
			break
		}
	}
}

////////////////////////////////////////////////////////
//
//客户端启动函数
//参数
//      远程ip地址和端口 tcpaddr
//
////////////////////////////////////////////////////////

func StartClient(tcpaddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)
	checkError(err, "ResolveTCPAddr")
	conn, err := net.DialTCP("tcp", nil, tcpAddr) // 服务端是 Listen and for{ Accept() } 客户端是 DialTCP
	// conn, err := net.DialTimeout("tcp", tcpAddr, 2 * time.Second)
	checkError(err, "DialTCP")
	//启动客户端发送线程
	go chatSend(conn)
	defer conn.Close()
	// 开始客户端轮训
	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf) // 循环（阻塞）监听从conn连接中读取服务端发送过来的数据
		if checkError(err, "Connection") == false {
			conn.Close()
			fmt.Println("Server is dead ...ByeBye")
			os.Exit(0)
		}
		fmt.Println("client get message: ", string(buf[0:length]))
	}
}

////////////////////////////////////////////////////////
//
//主程序
//
//参数说明：
//  启动服务器端：  Chat server [port]             eg: Chat server 9090
//  启动客户端：    Chat client [Server Ip Addr]:[Server Port]    eg: Chat client 192.168.0.74:9090
//
////////////////////////////////////////////////////////

// server: ./example1 server 9090
// client: ./example1 client localhost:9090

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Wrong pare")
		os.Exit(0)
	}
	if os.Args[1] == "server" && len(os.Args) == 3 {

		StartServer(os.Args[2])
	}
	if os.Args[1] == "client" && len(os.Args) == 3 {

		StartClient(os.Args[2])
	}
}
