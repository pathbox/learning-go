package main

import (
	"fmt"
	"net"
)

func main() {
	ip := GetIP()
	fmt.Println("The server IP: ", ip)
}

// get inet addr ip 内网ip地址
func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("ipnet IP: ", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
