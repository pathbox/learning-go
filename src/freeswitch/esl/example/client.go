package example

import (
	"flag"
	"fmt"
	"strings"

	. "github.com/0x19/goesl"
)

var (
	fshost   = flag.String("fshost", "localhost", "Freeswitch hostname. Default: localhost")
	fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
)

func main() {

	// Boost it as much as it can go ...
	// We don't need this since Go 1.5
	// runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := NewClient(*fshost, *fsport, *password, *timeout)

	if err != nil {
		Error("Error while creating new client: %s", err)
		return
	}

	// Apparently all is good... Let us now handle connection :)
	// We don't want this to be inside of new connection as who knows where it my lead us.
	// Remember that this is crutial part in handling incoming messages. This is a must!
	go client.Handle()

	client.Send("events json ALL")

	client.BgApi(fmt.Sprintf("originate %s %s", "sofia/internal/1001@127.0.0.1", "&socket(192.168.1.2:8084 async full)"))

	for {
		msg, err := client.ReadMessage()

		if err != nil {

			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				Error("Error while reading Freeswitch message: %s", err)
			}

			break
		}

		fmt.Printf("Got new message from Freeswitch response: %s\n", msg)
	}
}
