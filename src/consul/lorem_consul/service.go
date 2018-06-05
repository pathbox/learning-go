package lorem_consul

import (
	golorem "github.com/drhodes/golorem"
)

// Define service interface
type Service interface {
	// generate a word with at least min letters and at most max letters.
	Word(min, max int) string

	// generate a sentence with at least min words and at most max words.
	Sentence(min, max int) string

	// generate a paragraph with at least min sentences and at most max sentences.
	Paragraph(min, max int) string

	// health check
	HealthCheck() bool
}

// Implement service with empty struct
type LoremService struct {
}

// create type that return function.
// this will be needed in main.go
type ServiceMiddleware func(Service) Service

// Implement service functions
func (LoremService) Word(min, max int) string {
	return golorem.Word(min, max)
}

func (LoremService) Sentence(min, max int) string {
	return golorem.Sentence(min, max)
}

func (LoremService) Paragraph(min, max int) string {
	return golorem.Paragraph(min, max)
}

func (LoremService) HealthCheck() bool {
	//to make the health check always return true is a bad idea
	//however, I did this on purpose to make the sample run easier.
	//Ideally, it should check things such as amount of free disk space or free memory
	return true
}
