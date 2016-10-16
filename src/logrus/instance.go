package main

import (
	"os"
	"github.com/Sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instance.
var log = logrus.New()

func main() {
	log.Out = os.Stderr
	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size": 10,
		}).Info("A group of walrus emerges from the ocean")
}

