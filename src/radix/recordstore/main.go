package main

import (
	"./models"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/album", showAlbum)
	http.HandleFunc("like", addLike)
	http.handleFunc("/popular", listPopular)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func listPopular(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	abs, err := models.FindTopThree()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for i, ab := range abs {
		fmt.Fprintf(w, "%d) %s by %s: £%.2f [%d likes] \n", i+1, ab.Title, ab.Artist, ab.Price, ab.Likes)
	}
}

func addLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
		return
	}
	id := r.PostFormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// Call the IncrementLikes() function passing in the user-provided id. If
	// there's no album found with that id, return a 404 Not Found response.
	// In the event of any other errors, return a 500 Internal Server Error
	// response.
	err := models.IncrementLikes(id)
	if err == models.ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	// Redirect the client to the GET /ablum route, so they can see the
	// impact their like has had.
	http.Redirect(w, r, "/album?id="+id, 303)
}

func showAlbum(w http.ResponseWriter, r *http.Request) {
	// Unless the request is using the GET method, return a 405 'Method Not
	// Allowed' response.
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retrieve the id from the request URL query string. If there is no id
	// key in the query string then Get() will return an empty string. We
	// check for this, returning a 400 Bad Request response if it's missing.
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Call the FindAlbum() function passing in the user-provided id. If
	// there's no matching album found, return a 404 Not Found response. In
	// the event of any other errors, return a 500 Internal Server Error
	// response.
	bk, err := models.FindAlbum(id)
	if err == models.ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Write the album details as plain text to the client.
	fmt.Fprintf(w, "%s by %s: £%.2f [%d likes] \n", bk.Title, bk.Artist, bk.Price, bk.Likes)
}
