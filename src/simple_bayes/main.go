package main

import (
	"fmt"

	"./bayes"
)

func main() {
	var classifier bayes.Classifier
	classifier.Train([]string{"cat", "meow", "cat"}, "cat")
	classifier.Train([]string{"purr", "cat", "cat"}, "cat")
	classifier.Train([]string{"cat", "whisker", "meow"}, "cat")
	classifier.Train([]string{"dog", "bark", "cat"}, "dog")

	class, score := classifier.Classify([]string{"cat", "cat", "cat", "dog", "bark"})

	fmt.Println(class, score) // => cat 0.0002133333333333334
}
