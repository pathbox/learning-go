package main

import (
	"./controllers"
	"log"
	"net/http"
)

func main() {

	// myhandlerFunc := http.HandlerFunc(addValueToCmap)

	http.HandleFunc("/cmap", controllers.AddValueToCmap)
	log.Fatal(http.ListenAndServe(":9090", nil))

}
