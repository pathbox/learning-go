package main

import "fmt"

type Player interface {
	Jump()
	Run()
}

type Ao struct {
	Name string
}

type Bo struct {
	Name string
}

func (a *Ao) Jump() {
	fmt.Println(a.Name + " Jump")
}
func (b *Bo) Jump() {
	fmt.Println(b.Name + " Jump")
}

func (a *Ao) Run() {
	fmt.Println(a.Name + " Run")
}
func (b *Bo) Run() {
	fmt.Println(b.Name + " Run")
}

func main() {
	m := make(map[string]Player)

	a := &Ao{Name: "AAA"}
	b := &Bo{Name: "BBB"}

	m[a.Name] = a
	m[b.Name] = b

	for _, v := range m {
		v.Jump()
		v.Run()
	}
}
