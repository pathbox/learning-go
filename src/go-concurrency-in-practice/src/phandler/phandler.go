package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name    string
	Age     uint8
	Address Addr
}

type Addr struct {
	city     string
	district string
}

// 定义一个接口类型
type PersonHandler interface {
	Batch(origins <-chan Person) <-chan Person
	Handle(orig *Person)
}

type PersonHandlerImpl struct{}

func (handler PersonHandlerImpl) Batch(origs <-chan Person) <-chan Person {
	dests := make(chan Person, 100)
	go func() {
		for p := range origs {
			handler.Handle(&p) // 简单的修改Address.district 逻辑
			dests <- p
		}
		fmt.Println("All the information has been handled.")
		close(dests) // 这行注释会报 all goroutines are asleep - deadlock!
	}()
	return dests // 以 <-chan 返回chan
}

func (handler PersonHandlerImpl) Handle(orig *Person) {
	if orig.Address.district == "Haidian" {
		orig.Address.district = "Shijingshan"
	}
}

var personTotal = 200
var persons []Person = make([]Person, personTotal)
var personCount int

func init() {
	for i := 0; i < 200; i++ { // 初始化 200条数据
		name := fmt.Sprintf("%s%d", "p", i)
		p := Person{name, 32, Addr{"Beijing", "Haidian"}}
		persons[i] = p
	}
}

func main() {
	handler := getPersonHandler() // 返回PersonHandlerImpl
	origs := make(chan Person, 100)
	dests := handler.Batch(origs) // 2处理数据 从origs获得数据，进行处理
	fetchPerson(origs)            // 1创造数据 放入origs
	sign := savePerson(dests)     // 3保存数据 保存处理dests
	<-sign                        //阻塞等待savePerson的子goroutine完成
}

func getPersonHandler() PersonHandler {
	return PersonHandlerImpl{}
}

func savePerson(dest <-chan Person) <-chan byte { // chan<- Person 或 chan Person 都不行（invalid operation: <-dest (receive from send-only type chan<- Person）
	sign := make(chan byte, 1)
	go func() {
		for {
			p, ok := <-dest
			if !ok {
				fmt.Println("All the information has been saved")
				sign <- 0
				break
			}
			savePerson1(p)
		}
	}()
	return sign
}

func fetchPerson(origs chan<- Person) {
	origsCap := cap(origs)
	buffered := origsCap > 0
	goTicketTotal := origsCap / 2
	goTicket := initGoTicket(goTicketTotal)
	go func() {
		for {
			p, ok := fecthPerson1()
			if !ok {
				for {
					if !buffered || len(goTicket) == goTicketTotal {
						break
					}
					time.Sleep(time.Nanosecond)
				}
				fmt.Println("All the information has been fetched.")
				close(origs)
				break
			}
			if buffered {
				<-goTicket
				go func() {
					origs <- p
					goTicket <- 1
				}()
			} else {
				origs <- p
			}
		}
	}()
}

func initGoTicket(total int) chan byte {
	var goTicket chan byte
	if total == 0 {
		return goTicket
	}
	goTicket = make(chan byte, total)
	for i := 0; i < total; i++ {
		goTicket <- 1
	}
	return goTicket
}

func fecthPerson1() (Person, bool) {
	if personCount < personTotal {
		p := persons[personCount]
		personCount++
		return p, true
	}
	return Person{}, false
}

func savePerson1(p Person) bool {
	return true
}

/* All the information has been fetched.
All the information has been handled.
All the information has been saved
*/
