type appHandler func(http.ResponseWriter, *http.Request)

// Our appHandler type will now satisify http.Handler
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if status, err := fn(w, r); err != nil {
    // We could also log our errors centrally:
        // i.e. log.Printf("HTTP %d: %v", err)
    switch status {
      // We can have cases as granular as we like, if we wanted to
        // return custom errors for specific status codes.
    case http.StatusNotFound:
      notFound(w, r)
    case http.StatusInternalServerError:
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    default:
      // Catch any other errors we haven't explicitly handled
      http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    }
  }
}
func myHandler(w http.ResponseWriter, r *http.Request) (int, error) {
    session, err := store.Get(r, "myapp")
    if err != nil {
      return http.StatusInternalServerError, err
    }

    post := Post{ID: id}
    exists, err := db.GetPost(&post)
    if err != nil {
      return http.StatusInternalServerError, err
    }

    // We can shortcut this: since renderTemplate returns `error`,
    // our ServeHTTP method will return a HTTP 500 instead and won't
    // attempt to write a broken template out with a HTTP 200 status.
    // (see the postscript for how renderTemplate is implemented)
    // If it doesn't return an error, things will go as planned.
    return http.StatusOK, renderTemplate(w, "post.tmpl", data)
  }

func main() {
  // Cast myHandler to an appHandler
  http.Handle("/", appHandler(myHandler))
  http.ListenAndServe(":9090", nil)
}