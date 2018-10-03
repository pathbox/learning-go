package main

import (
	"fmt"

	"net"

	"os"

	"encoding/gob"
)

type Person struct {
	Name

	Email []Email
}

type Name struct {
	Family string

	Personal string
}

type Email struct {
	Kind string

	Address string
}

func (p Person) String() string {

	s := p.Name.Personal + " " + p.Name.Family

	for _, v := range p.Email {

		s += "\n" + v.Kind + ": " + v.Address

	}

	return s

}

func main() {

	service := "0.0.0.0:9091"

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)

	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)

	checkError(err)

	for {

		conn, err := listener.Accept()

		if err != nil {

			continue

		}

		encoder := gob.NewEncoder(conn)

		decoder := gob.NewDecoder(conn)

		for n := 0; n < 10; n++ {

			var person Person

			decoder.Decode(&person)

			fmt.Println(person.String())

			encoder.Encode(person)

		}

		conn.Close() // we're finished

	}

}

func checkError(err error) {

	if err != nil {

		fmt.Println("Fatal error ", err.Error())

		os.Exit(1)

	}

}
