package abstract_factory

import "testing"

func Test(t *testing.T) {
	t.Run("abstract_factory: ", ProduceFruitAndEat)
}

func ProduceFruitAndEat(t *testing.T) {
	var factory FruitFactory
	var apple Fruit

	factory = &AppleFactory{}

	apple = factory.CreateFruit()

	apple.Eat()
}
