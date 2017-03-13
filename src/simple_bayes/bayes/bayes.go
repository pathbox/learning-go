package bayes

// Classifier implements the Naive Bayes Classifier.
type Classifier struct {
	classPerWordCount map[string]int
	classPerDocument  map[string]int
	wordClassCount    map[string]map[string]int
	documentsCount    int
}

// Train the classifier with text and its class.
func (c *Classifier) Train(words []string, class string) {
	c.documentsCount++
	c.incrementDocumentPerClass(class)

	for _, word := range words {
		c.incrementWordClass(word, class)
		c.incrementClassPerWord(class)
	}
}

// Classify the given text.
func (c *Classifier) Classify(words []string) (string, float64) {
	var score float64
	var prediction string

	for _, class := range c.classes() {
		var probability = c.probability(words, class)

		if score < probability {
			score = probability
			prediction = class
		}
	}

	return prediction, score
}

// counting
func (c *Classifier) incrementWordClass(word string, class string) {
	if c.wordClassCount == nil {
		c.wordClassCount = make(map[string]map[string]int)
	}
	if c.wordClassCount[word] == nil {
		c.wordClassCount[word] = make(map[string]int)
	}
	c.wordClassCount[word][class]++
}

func (c *Classifier) incrementDocumentPerClass(class string) {
	if c.classPerDocument == nil {
		c.classPerDocument = make(map[string]int)
	}
	c.classPerDocument[class]++
}

func (c *Classifier) incrementClassPerWord(class string) {
	if c.classPerWordCount == nil {
		c.classPerWordCount = make(map[string]int)
	}
	c.classPerWordCount[class]++
}

func (c *Classifier) classProbability(class string) float64 {
	return float64(c.classPerDocument[class]) / float64(c.documentsCount)
}

func (c *Classifier) wordGivenClassProbability(word string, class string) float64 {
	return float64(c.wordClassCount[word][class]+1) / float64(c.classPerWordCount[class]+c.vocabularySize())
}

func (c *Classifier) probability(words []string, class string) float64 {
	var result = c.classProbability(class)

	for _, word := range words {
		result *= c.wordGivenClassProbability(word, class)
	}

	return result
}

// helpers
func (c *Classifier) classes() []string {
	classes := make([]string, len(c.classPerDocument))

	i := 0
	for c := range c.classPerDocument {
		classes[i] = c
		i++
	}

	return classes
}

func (c *Classifier) vocabularySize() int {
	return len(c.wordClassCount)
}
