package abstract_factory // 抽象工程模式

import "fmt"

type FruitFactory interface {
	CreateFruit() Fruit
}

// ApppleFactory 苹果工厂，实现FruitFactory接口
type AppleFactory struct{}

func (apppleFactory AppleFactory) CreateFruit() Fruit {
	return &Apple()
}

// Fruit 水果接口
type Fruit interface {
	Eat()
}

// Apple 苹果，实现Fruit接口
type Apple struct{}

func (apple Apple) Eat() {
	fmt.Printl("Eat Apple")
}
