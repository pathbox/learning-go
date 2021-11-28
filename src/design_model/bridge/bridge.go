package bridge

import "fmt"

type ICoffee interface {
	OrderCoffee()
}

type LargeCoffee struct {
	ICoffeeAddtion
}

type MediumCoffee struct {
	ICoffeeAddtion
}

type SmallCoffee struct {
	ICoffeeAddtion
}

func (lc LargeCoffee) OrderCoffee() {
	fmt.Println("big one")
	lc.AddSomething()
}

func (mc MediumCoffee) OrderCoffee() {
	fmt.Println("medium one")
	mc.AddSomething()
}

func (sc SmallCoffee) OrderCoffee() {
	fmt.Println("small one")
	sc.AddSomething()
}


type CoffeeCupType uint8

const (
	CoffeeCupTypeLarge = iota
	CoffeeCupTypeMedium = iota
	CoffeeCupTypeSmall = iota
)

var CoffeeFuncMap = map[CoffeeCupType]func(coffeeAddtion ICoffeeAddtion) ICoffee {
	CoffeeCupTypeLarge: NewLargeCoffee,
	CoffeeCupTypeMedium: NewMediumCoffee,
	CoffeeCupTypeSmall: NewSmallCoffee,
}

func NewCoffee(cupType CoffeeCupType, coffeeAddtion ICoffeeAddtion) ICoffee {
	if handler, ok := CoffeeFuncMap[cupType]; ok {
		return handler(coffeeAddtion)
	}
	return nil
}

func NewLargeCoffee(coffeeAddtion ICoffeeAddtion) ICoffee {
	return &LargeCoffee{coffeeAddtion}
}

func NewMediumCoffee(coffeeAddtion ICoffeeAddtion) ICoffee {
	return &MediumCoffee{coffeeAddtion}
}

func NewSmallCoffee(coffeeAddtion ICoffeeAddtion) ICoffee {
	return &SmallCoffee{coffeeAddtion}
}

type Milk struct{}

type Sugar struct{}

func (milk Milk) AddSomething() {
	fmt.Println("add milk")
}

func (sugar Sugar) AddSomething() {
	fmt.Println("add sugar")
}

type coffeeAddtionType uint8

const {
	CoffeeAddtionTypeMilk = iota
	CoffeeAddtionTypeSugar = iota
}
