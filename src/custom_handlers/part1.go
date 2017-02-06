// We create our own handler type that satisfies http.Handler (read: it has a  ServeHTTP(http.ResponseWriter, *http.Request) method), which allows it to remain compatible with net/http, generic HTTP middleware packages like nosurf, and routers/frameworks like gorilla/mux or Goji.

func myHandler(w http.ResponseWriter, r *http.Request) {
  session, err := store.Get(r, "myapp")
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return // Forget to return, and the handler will continue on
  }

  id := // get id from URL param; strconv.AtoI it; making sure to return on those errors too...
  post := Post{ID: id}
  exists, err := db.GetPost(&post)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  if !exists {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        return // ... and again.
    }

    err = renderTemplate(w, "post.tmpl", post)
    if err != nil {
        // Yep, here too...
    }
}