package main

import log "github.com/Sirupsen/logrus"

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.WithFields(log.Fields{
		"animal": "walrus",
		"number": 1,
		"size": 10,
		"slice": []string{"a", "b"},
		}).Info("A walrus appears")
	log.Info("Another info")
}
