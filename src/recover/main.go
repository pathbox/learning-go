package main

import (
	"log"
	"net/http"

	"github.com/dre1080/recover"
)

var myPanicHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	panic("You should not have a handler the just panics")
})

func main() {
	recovery := recover.New(&recover.Options{
		Log: log.Println,
	})
	// recoveryWithDefaults := recovery.New(nil)

	app := recovery(myPanicHandler)
	log.Println("Listening at :9090")
	log.Println(http.ListenAndServe(":9090", app))
}
