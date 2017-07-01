func server(ctx context.Context, wg *sync.WaitGroup) {
  defer wg.Done()

  mux := http.NewServeMux()
  mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    fmt.Println("server: received request")
    time.Sleep(3 * time.Second)
    io.WriteString(w, "Finished!\n")
    fmt.Println("server: request finished")
    }))

  srv := &http.Server{Addr: ":8080", Handler: mux}
  go func(){
    if err := srv.ListenAndServer(); err != nil {
      fmt.Printf("Listen : %s\n", err)
    }
  }()

  <-ctx.Done()
  fmt.Println("server: caller has told us to stop")

    // shut down gracefully, but wait no longer than 5 seconds before halting
  shutdownCtx, cancel := context.WithTimeout(context.Background(). 5*time.Second)
  defer cancel()

  srv.Shutdown(shutdownCtx)

  fmt.Println("server gracefully stopped")
}