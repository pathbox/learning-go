package main

import (
	"os"
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stderr)

	log.SetLevel(log.WarnLevel)
}

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
		"size": 10,
		}).Info("A group of walrus emerges form the ocean")
	log.WithFields(log.Fields{
		"omg": true,
		"number": 122,
		}).Warn("The gtoup's numbre increased tremendously!")

	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other": "I also should be logged always",
		})

	contextLogger.Warn("I will be logged with common and other field" )
	contextLogger.Warn("Me too")

	log.WithFields(log.Fields{
		"omg": true,
		"number": 100,
		}).Fatal("The ice breaks!")
}
