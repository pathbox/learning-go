package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tidwall/evio"
)

func main() {
	var pduplicates, puniques int
	var ticks int
	var logidx int
	var requests int
	var events evio.Events
	streams := make(map[int]*evio.InputStream)
	numbers := make(map[int]int)

	for i := 0; ; i++ {
		if err := os.Remove(fmt.Sprintf("data.%d.log", i)); err != nil {
			if os.IsNotExist(err) {
				break
			}
			log.Fatal(err)
		}
	}
	f, err := os.Create("data.0.log")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	events.Opened = func(id int, _ evio.Info) (_ []byte, _ evio.Options, action evio.Action) {
		if len(streams) == 6 {
			action = evio.Close
		} else {
			streams[id] = new(evio.InputStream)
		}
		return
	}
	events.Closed = func(id int, _ error) (_ evio.Action) {
		delete(streams, id)
		return
	}
	events.Data = func(id int, in []byte) (out []byte, action evio.Action) {
		s := streams[id]
		data := s.Begin(in)
	nextLine:
		if len(data) >= 8 && data[0] == 's' && string(data[:8]) == "shutdown" {
			log.Printf("shutdown received")
			return nil, evio.Shutdown
		}
		var number int
		for i, j := 0, 0; i < len(data); i++ {
			switch {
			default:
				number, j = number*10+int(data[i]-'0'), j+1
			case data[i] == '\r' && i < len(data)-1 && data[i+1] == '\n':
				i++
				fallthrough
			case data[i] == '\n':
				n := len(numbers)
				numbers[number]++
				if n == len(numbers) {
					pduplicates++
				} else {
					puniques++
					if _, err := f.Write(append(append([]byte{}, data[:j]...), '\n')); err != nil {
						log.Fatal(err)
					}
				}
				requests++
				data = data[i+1:]
				goto nextLine
			case data[i] < '0' || data[i] > '9': // byte is not number char
				log.Printf("malformed data")
				return nil, evio.Close
			}
		}
		s.End(data)
		return
	}

	events.Serving = func(_ evio.Server) (_ evio.Action) {
		log.Printf("server started on port 5055")
		return
	}

	events.Tick = func() (delay time.Duration, action evio.Action) {
		if ticks > 0 {
			if ticks%5 == 4 {
				log.Printf("period: (unique: %d, duplicates: %d), totals: (uniques: %d, requests: %d)", puniques, pduplicates, len(numbers), requests)
				puniques, pduplicates = 0, 0
			}
			if ticks%10 == 9 {
				log.Printf("rotating log: data.%d.log", logidx)
				if err := f.Close(); err != nil {
					log.Fatal(err)
				}
				logidx++
				if f, err = os.Create(fmt.Sprintf("data.%d.log", logidx)); err != nil {
					log.Fatal(err)
				}
			}
		}
		ticks++
		delay = time.Second
		return
	}
	fmt.Println("evio server Starting...")
	if err := evio.Serve(events, "tcp://:5055"); err != nil {
		log.Fatal(err)
	}
}
