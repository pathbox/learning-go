package controllers

import (
	"fmt"
	"github.com/orcaman/concurrent-map"
	"net/http"
	"strconv"
	"time"
)

var Cmap = cmap.New()

func AddValueToCmap(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Unix()
	uuid := fmt.Sprintf("%v", timestamp)
	fmt.Println(uuid)

	Cmap.Set(uuid, "true")
	count := Cmap.Count()

	show_count := strconv.Itoa(count)
	w.Write([]byte(show_count))
}
