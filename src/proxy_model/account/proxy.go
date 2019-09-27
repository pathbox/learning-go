package proxy

import (
    "interceptor/account"
    "fmt"
)

type Proxy struct {
  Account account.Account
}

func (p *Proxy) Query(id string) int {
	fmt.Println("Proxy.Query begin")
	value := p.Account.Query(id)
	fmt.Println("Proxy.Query end")
	return value
}

func (p *Proxy) Update(id string, value int) {
	fmt.Println("Proxy.Update begin")
	p.Account.Update(id, value)
	fmt.Println("Proxy.Update end")
}