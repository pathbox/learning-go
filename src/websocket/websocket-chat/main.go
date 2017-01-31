package main

mport (
  "log"
  "net/http"

  "./chat"
)

func main() {
  log.SetFlags(log.Lshortfile)

  server := chat.NewServer("/entry")
  go server.Listen()

  http.Handle("/", http.FileServer(http.Dir("webroot")))

  log.Fatal(http.ListenAndServe(":9090", nil))
}

