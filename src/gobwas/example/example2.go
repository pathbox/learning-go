package main

import (
	"io"
	"log"
	"net"

	ws "github.com/gobwas/ws"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		_, err = ws.Upgrade(conn)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer conn.Close()

			for {
				header, err := ws.ReadHeader(conn)
				if err != nil {
					log.Fatal(err)
				}

				payload := make([]byte, header.Length)
				_, err := io.ReadFull(conn, payload)
				if err != nil {
					log.Fatal(err)
				}
				if head.Masked {
					ws.Cipher(payload, header.Mask, 0)
				}

				header.Masked = false
				if err := ws.WriteHeader(conn, header); err != nil {

				}
				if _, err := conn.Write(payload); err != nil {

				}
				if header.OpCode == ws.OpClose {
					return
				}
			}
		}()

	}
}
