import "context"

/* src/net/tcpsock.go
ResolveTCPAddr()接收两个字符串参数。

network: 必须是TCP网络名，比如tcp, tcp4, tcp6。
address: TCP地址字符串，如果它不是字面量的IP地址或者端口号不是字面量的端口号， ResolveTCPAddr会将传入的地址解决成TCP终端的地址。否则传入一对字面量IP地址和端口数字作为地址。address参数可以使用host名称，但是不推荐这样做，因为它最多会返回host名字的一个IP地址。
ResolveTCPAddr()接收的代表TCP地址的字符串(例如localhost:80, 127.0.0.1:80, 或[::1]:80, 都是代表本机的80端口), 返回(net.TCPAddr指针, nil)(如果字符串不能被解析成有效的TCP地址会返回(nil, error))。
*/

func ResolveTCPAddr(network, address string) (*TCPAddr, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
	case "":
		network = "tcp"
	default:
		return nil, UnknownNetworkError(network)
	}
	addrs, err := DefaultResolver.internetAddrList(context.Background(), network, address)
	if err != nil {
		return nil, err
	}
	return addrs.forResolve(network, address).(*TCPAddr), nil
}