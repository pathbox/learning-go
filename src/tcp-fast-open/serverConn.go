// Interfaces for the server's establish tcp connection with a client

package main

import (
	"log"
	"syscall"
)

type TFOServerConn struct {
	sockaddr *syscall.SockaddrInet4
	fd       int
}

// Read the data from the client and immediately close the connection
func (cxn *TFOServerConn) Handle() {
	defer cxn.Close()

	log.Printf("Server Conn: Connection received from remote addr: %v, remote port: %d\n",
		cxn.sockaddr.Addr, cxn.sockaddr.Port)

	// Create a small buffer to store data from client
	buf := make([]byte, 24)

	n, err := syscall.Read(cxn.fd, buf)
	if err != nil {
		log.Println("Failed to read() client:", err)
		return
	}

	// Do nothing in particular with the response, just print it
	log.Printf("Server Conn: Read %d bytes: %#v", n, string(buf[:n]))

	// The defer will close the connection now
}

func (cxn *TFOServerConn) Close() {
	// 关掉连接
	err := syscall.Shutdown(cxn.fd, syscall.SHUT_RDWR)
	if err != nil {
		log.Println("Failed to shutdown() connection:", err)
	}

	// Close the file descriptor
	err = syscall.Close(cxn.fd)
	if err != nil {
		log.Println("Failed to close() connection:", err)
	}

}
