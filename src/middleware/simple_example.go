type SingleHost struct {
  handler http.Handler
  allowedHost string
}

func NewSingleHost(handler http.Handler, allowedHost string) *SingleHost {
  return &SingleHost{handler: handler, allowedHost: allowedHost}
}

/* go source code
Now, for the actual logic. To implement http.Handler, our type only needs to have one method

type Handler interface {
  ServeHTTP(ResponseWriter, *Request)
}
*/

func (s *SingleHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  host := r.Host
  if host == s.allowedHost {
    s.handler.ServeHTTP(w, r)
  } else {
    w.WriteHeader(403)
  }
}

singleHosted = NewSingleHost(myHandler, "example.com")
http.ListenAndServe(":8080", singleHosted)


// middleware

func SingleHost(handler http.Handler, allowedHost string) http.Handler {
  ourFunc := func(w http.ResponseWriter, r *http.Request) {
    host := r.Host
    if host == allowedHost{
      handler.ServeHTTP(w, r)
    } else {
      w.WriteHeader(403)
    }
  }
  return http.HandlerFunc(ourFunc)
}


type AppendMiddleware struct{
  handler http.Handler
}

func (a *AppendMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  a.handler.ServeHTTP(w, r)
  w.Write([]byte("Middleware says hello"))
}

type ModifierMiddleware struct {
    handler http.Handler
}

func (m *ModifierMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rec := httptest.NewRecorder()
    // passing a ResponseRecorder instead of the original RW
    m.handler.ServeHTTP(rec, r)
    // after this finishes, we have the response recorded
    // and can modify it before copying it to the original RW

    // we copy the original headers first
    for k, v := range rec.Header() {
        w.Header()[k] = v
    }
    // and set an additional one
    w.Header().Set("X-We-Modified-This", "Yup")
    // only then the status code, as this call writes out the headers
    w.WriteHeader(418)

    // The body hasn't been written (to the real RW) yet,
    // so we can prepend some data.
    data := []byte("Middleware says hello again. ")

    // But the Content-Length might have been set already,
    // we should modify it by adding the length
    // of our own data.
    // Ignoring the error is fine here:
    // if Content-Length is empty or otherwise invalid,
    // Atoi() will return zero,
    // which is just what we'd want in that case.
    clen, _ := strconv.Atoi(r.Header.Get("Content-Length"))
    clen += len(data)
    r.Header.Set("Content-Length", strconv.Itoa(clen))

    // finally, write out our data
    w.Write(data)
    // then write out the original body
    w.Write(rec.Body.Bytes())
}

// sharing data with other handlers

type csrfContext struct {
  token string
  reason error
}

var (
  contextMap = make(map[*http.Request]*csrfContext)
  cmMutex = new(sync.RWMutex)
)

func Token(req *http.Request) string {
  cmMutex.RLock()
  defer cmMutex.RUlock()

  ctx, ok := contextMap[req]
  if !ok {
    return ""
  }
  return ctx.token
}

