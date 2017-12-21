package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/heroku/x/redis"
)

var (
	waitTimeout = time.Minute * 10
	log         = logrus.WithField("cmd", "go-websocket-chat-demo")
	rr          redisReceiver
	rw          redisWriter
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.WithField("REDIS_URL", redisURL).Fatal("$REDIS_URL must be set")
	}
	redisPool, err := redis.NewRedisPoolFromURL(redisURL)
	if err != nil {
		log.WithField("url", redisURL).Fatal("Unable to create Redis pool")
	}

	rr = newRedisReceiver(redisPool)
	rw = newRedisWriter(redisPool)

	go func() {
		for {
			waited, err := redis.WaitForAvailability(redisURL, waitTimeout, rr.wait)
			if !waited || err != nil {
				log.WithFields(logrus.Fields{"waitTimeout": waitTimeout, "err": err}).Fatal("Redis not available by timeout!")
			}
			rr.broadcast(availableMessage)
			err = rr.run()
			if err == nil {
				break
			}
			log.Error(err)
		}
	}()

	go func() {
		for {
			waited, err := redis.WaitForAvailability(redisURL, waitTimeout, nil)
			if !waited || err != nil {
				log.WithFields(logrus.Fields{"waitTimeout": waitTimeout, "err": err}).Fatal("Redis not available by timeout!")
			}
			err = rw.run()
			if err == nil {
				break
			}
			log.Error(err)
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleWebsocket)
	log.Println(http.ListenAndServe(":"+port, nil))
}
