package sshmux

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// DefaultInteractive is the default server selection prompt for users during
// session forward.
func DefaultInteractive(comm io.ReadWriter, session *Session) (string, error) {
	remotes := session.Remotes

	fmt.Fprintf(comm, "Welcome to sshmux, %s\r\n", session.Conn.User())
	for i, v := range remotes {
		fmt.Fprintf(comm, "    [%d] %s\r\n", i, v)
	}

	// Beware, nasty input parsing loop
loop:
	for {
		fmt.Fprintf(comm, "Please select remote server: ")
		var buf []byte
		b := make([]byte, 1)
		var (
			n   int
			err error
		)
		for {
			if err != nil {
				return "", err
			}
			n, err = comm.Read(b)
			if n == 1 {
				fmt.Fprintf(comm, "%s", b)
				switch b[0] {
				case '\r':
					fmt.Fprintf(comm, "\r\n")
					res, err := strconv.ParseInt(string(buf), 10, 64)
					if err != nil {
						fmt.Fprintf(comm, "input not a valid integer. Please try again\r\n")
						continue loop
					}
					if int(res) >= len(remotes) || res < 0 {
						fmt.Fprintf(comm, "No such server. Please try again\r\n")
						continue loop
					}

					return remotes[int(res)], nil
				case 0x03:
					fmt.Fprintf(comm, "\r\nGoodbye\r\n")
					return "", errors.New("user terminated session")
				}
				buf = append(buf, b[0])
			}
		}
	}
}
