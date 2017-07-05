package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		io.WriteString(w, "Finished!")
	}))

	srv := &http.Server{Addr: ":8081", Handler: mux}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")

}
