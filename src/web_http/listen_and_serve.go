func ListenAndServe(addr string, handler Handler) error {
  server := &Server{Addr: addr, Handler: handler}
  return server.ListenAndServe()
}


// 这个函数其实也是一层封装，创建了 Server 结构，并调用它的 ListenAndServe 方法，那我们就跟进去看看：

type Server struct {
  Addr string
  Handler Handler
  ......
}

func(srv *Server) ListenAndServe() error {
  addr := srv.Addr
  if addr == ""{
    addr = ":http"
  }
  In, err := net.Listen("tcp", addr)
  if err != nil {
    return err
  }
  return srv.Serve(tcpKeepAliveListener{In.(*net.TCPListener)})
}

// Server 保存了运行 HTTP 服务需要的参数，调用 net.Listen 监听在对应的 tcp 端口，tcpKeepAliveListener 设置了 TCP 的 KeepAlive 功能，最后调用 srv.Serve()方法开始真正的循环逻辑。我们再跟进去看看 Serve 方法：

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
  defer l.Close()
  var tempDelay time.Duration // how long to sleep on accept failure
    // 循环逻辑，接受请求并处理
  for {
         // 有新的连接
    rw, e := l.Accept()
    if e != nil {
      if ne, ok := e.(net.Error); ok && ne.Temporary() {
        if tempDelay == 0 {
          tempDelay = 5 * time.Millisecond
        } else {
          tempDelay *= 2
        }
        if max := 1 * time.Second; tempDelay > max {
          tempDelay = max
        }
        srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
        time.Sleep(tempDelay)
        continue
      }
      return e
    }
    tempDelay = 0
         // 创建 Conn 连接
    c, err := srv.newConn(rw)
    if err != nil {
      continue
    }
    c.setState(c.rwc, StateNew) // before Serve can return
         // 启动新的 goroutine 进行处理
    go c.serve()
  }
}

// 最上面的注释也说明了这个方法的主要功能：

// 接受 Listener l 传递过来的请求
// 为每个请求创建 goroutine 进行后台处理
// goroutine 会读取请求，调用 srv.Handler

func (c *conn) serve() {
  origConn := c.rwc // copy it before it's set nil on Close or Hijack

    ...

  for {
    w, err := c.readRequest()
    if c.lr.N != c.server.initialLimitedReaderSize() {
      // If we read any bytes off the wire, we're active.
      c.setState(c.rwc, StateActive)
    }

         ...

    // HTTP cannot have multiple simultaneous active requests.[*]
    // Until the server replies to this request, it can't read another,
    // so we might as well run the handler in this goroutine.
    // [*] Not strictly true: HTTP pipelining.  We could let them all process
    // in parallel even if their responses need to be serialized.
    serverHandler{c.server}.ServeHTTP(w, w.req)

    w.finishRequest()
    if w.closeAfterReply {
      if w.requestBodyLimitHit {
        c.closeWriteAndWait()
      }
      break
    }
    c.setState(c.rwc, StateIdle)
  }
}

func (sh serverHadnler) ServeHTTP(rw ResponseWriter, req *Request) {
  handler := sh.srv.Handler
  if handler == nil {
    handler = DefaultServeMux
  }
  if req.RequestURI == "*" && req.Method == "OPTIONS" {
    handler = globalOptionsHandler{}
  }
  handler.ServeHTTP(rw, req)
}

// 如果没有 handler 为空，就会使用它。handler.ServeHTTP(rw, req)，Handler 接口都要实现 ServeHTTP 这个方法，因为这里就要被调用啦。

// 我们已经知道，ServeMux 会以某种方式保存 URL 和 Handlers 的对应关系

type ServeMux struct {
  mu sync.RWMutex
  m map[string]muxEntry
  hosts bool
}

type muxEntry struct {
  expliclt bool
  h Handler
  pattern string
}

// 数据结构也比较直观，和我们想象的差不多，路由信息保存在字典中，接下来就看看几个重要的操作：路由信息是怎么注册的？ServeHTTP 方法到底是怎么做的？路由查找过程是怎样的？

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *ServeMux) Handle(pattern string, handler Handler) {
  mux.mu.Lock()
  defer mux.mu.Unlock()

    // 边界情况处理
  if pattern == "" {
    panic("http: invalid pattern " + pattern)
  }
  if handler == nil {
    panic("http: nil handler")
  }
  if mux.m[pattern].explicit {
    panic("http: multiple registrations for " + pattern)
  }

    // 创建 `muxEntry` 并添加到路由字典中
  mux.m[pattern] = muxEntry{explicit: true, h: handler, pattern: pattern}

  if pattern[0] != '/' {
    mux.hosts = true
  }

    // 这是一个很有用的小技巧，如果注册了 `/tree/`， `serveMux` 会自动添加一个 `/tree` 的路径并重定向到 `/tree/`。当然这个 `/tree` 路径会被用户显示的路由信息覆盖。
  // Helpful behavior:
  // If pattern is /tree/, insert an implicit permanent redirect for /tree.
  // It can be overridden by an explicit registration.
  n := len(pattern)
  if n > 0 && pattern[n-1] == '/' && !mux.m[pattern[0:n-1]].explicit {
    // If pattern contains a host name, strip it and use remaining
    // path for redirect.
    path := pattern
    if pattern[0] != '/' {
      // In pattern, at least the last character is a '/', so
      // strings.Index can't be -1.
      path = pattern[strings.Index(pattern, "/"):]
    }
    mux.m[pattern[0:n-1]] = muxEntry{h: RedirectHandler(path, StatusMovedPermanently), pattern: pattern}
  }
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
  if r.RequestURI == "*" {
    if r.ProtoAtLeast(1, 1) {
      w.Header().Set("Connection", "close")
    }
    w.WriteHeader(StatusBadRequest)
    return
  }
  h, _ := mux.Handler(r)
  h.ServeHTTP(w, r)
}
// ServeHTTP 也只是通过 mux.Handler(r) 找到请求对应的 handler，调用它的 ServeHTTP 方法，代码比较简单我们就显示了，它最终会调用 mux.match() 方法

/ Does path match pattern?
func pathMatch(pattern, path string) bool {
  if len(pattern) == 0 {
    // should not happen
    return false
  }
  n := len(pattern)
  if pattern[n-1] != '/' {
    return pattern == path
  }
    // 匹配的逻辑很简单，path 前面的字符和 pattern 一样就是匹配
  return len(path) >= n && path[0:n] == pattern
}

// Find a handler on a handler map given a path string
// Most-specific (longest) pattern wins
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
  var n = 0
  for k, v := range mux.m {
    if !pathMatch(k, path) {
      continue
    }
         // 最长匹配的逻辑在这里
    if h == nil || len(k) > n {
      n = len(k)
      h = v.h
      pattern = v.pattern
    }
  }
  return
}
// match 会遍历路由信息字典，找到所有匹配该路径最长的那个

// Handle registers the handler for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }

// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
  DefaultServeMux.HandleFunc(pattern, handler)
}

// 原来是直接通过 DefaultServeMux 调用对应的方法，到这里上面的一切都串起来了！

// Request 就是封装好的客户端请求，包括 URL，method，header 等等所有信息，以及一些方便使用的方法：

// A Request represents an HTTP request received by a server
// or to be sent by a client.
//
// The field semantics differ slightly between client and server
// usage. In addition to the notes on the fields below, see the
// documentation for Request.Write and RoundTripper.
type Request struct {
  // Method specifies the HTTP method (GET, POST, PUT, etc.).
  // For client requests an empty string means GET.
  Method string

  // URL specifies either the URI being requested (for server
  // requests) or the URL to access (for client requests).
  //
  // For server requests the URL is parsed from the URI
  // supplied on the Request-Line as stored in RequestURI.  For
  // most requests, fields other than Path and RawQuery will be
  // empty. (See RFC 2616, Section 5.1.2)
  //
  // For client requests, the URL's Host specifies the server to
  // connect to, while the Request's Host field optionally
  // specifies the Host header value to send in the HTTP
  // request.
  URL *url.URL

  // The protocol version for incoming requests.
  // Client requests always use HTTP/1.1.
  Proto      string // "HTTP/1.0"
  ProtoMajor int    // 1
  ProtoMinor int    // 0

  // A header maps request lines to their values.
  // If the header says
  //
  //  accept-encoding: gzip, deflate
  //  Accept-Language: en-us
  //  Connection: keep-alive
  //
  // then
  //
  //  Header = map[string][]string{
  //    "Accept-Encoding": {"gzip, deflate"},
  //    "Accept-Language": {"en-us"},
  //    "Connection": {"keep-alive"},
  //  }
  //
  // HTTP defines that header names are case-insensitive.
  // The request parser implements this by canonicalizing the
  // name, making the first character and any characters
  // following a hyphen uppercase and the rest lowercase.
  //
  // For client requests certain headers are automatically
  // added and may override values in Header.
  //
  // See the documentation for the Request.Write method.
  Header Header

  // Body is the request's body.
  //
  // For client requests a nil body means the request has no
  // body, such as a GET request. The HTTP Client's Transport
  // is responsible for calling the Close method.
  //
  // For server requests the Request Body is always non-nil
  // but will return EOF immediately when no body is present.
  // The Server will close the request body. The ServeHTTP
  // Handler does not need to.
  Body io.ReadCloser

  // ContentLength records the length of the associated content.
  // The value -1 indicates that the length is unknown.
  // Values >= 0 indicate that the given number of bytes may
  // be read from Body.
  // For client requests, a value of 0 means unknown if Body is not nil.
  ContentLength int64

  // TransferEncoding lists the transfer encodings from outermost to
  // innermost. An empty list denotes the "identity" encoding.
  // TransferEncoding can usually be ignored; chunked encoding is
  // automatically added and removed as necessary when sending and
  // receiving requests.
  TransferEncoding []string

  // Close indicates whether to close the connection after
  // replying to this request (for servers) or after sending
  // the request (for clients).
  Close bool

  // For server requests Host specifies the host on which the
  // URL is sought. Per RFC 2616, this is either the value of
  // the "Host" header or the host name given in the URL itself.
  // It may be of the form "host:port".
  //
  // For client requests Host optionally overrides the Host
  // header to send. If empty, the Request.Write method uses
  // the value of URL.Host.
  Host string

  // Form contains the parsed form data, including both the URL
  // field's query parameters and the POST or PUT form data.
  // This field is only available after ParseForm is called.
  // The HTTP client ignores Form and uses Body instead.
  Form url.Values

  // PostForm contains the parsed form data from POST or PUT
  // body parameters.
  // This field is only available after ParseForm is called.
  // The HTTP client ignores PostForm and uses Body instead.
  PostForm url.Values

  // MultipartForm is the parsed multipart form, including file uploads.
  // This field is only available after ParseMultipartForm is called.
  // The HTTP client ignores MultipartForm and uses Body instead.
  MultipartForm *multipart.Form

  ...

  // RemoteAddr allows HTTP servers and other software to record
  // the network address that sent the request, usually for
  // logging. This field is not filled in by ReadRequest and
  // has no defined format. The HTTP server in this package
  // sets RemoteAddr to an "IP:port" address before invoking a
  // handler.
  // This field is ignored by the HTTP client.
  RemoteAddr string
    ...
}

// ResponseWriter

// ResponseWriter 是一个接口，定义了三个方法：

// Header()：返回一个 Header 对象，可以通过它的 Set() 方法设置头部，注意最终返回的头部信息可能和你写进去的不完全相同，因为后续处理还可能修改头部的值（比如设置 Content-Length、Content-type 等操作）
// Write()： 写 response 的主体部分，比如 html 或者 json 的内容就是放到这里的
// WriteHeader()：设置 status code，如果没有调用这个函数，默认设置为 http.StatusOK， 就是 200 状态码

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ResponseWriter interface {
  // Header returns the header map that will be sent by WriteHeader.
  // Changing the header after a call to WriteHeader (or Write) has
  // no effect.
  Header() Header

  // Write writes the data to the connection as part of an HTTP reply.
  // If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
  // before writing the data.  If the Header does not contain a
  // Content-Type line, Write adds a Content-Type set to the result of passing
  // the initial 512 bytes of written data to DetectContentType.
  Write([]byte) (int, error)

  // WriteHeader sends an HTTP response header with status code.
  // If WriteHeader is not called explicitly, the first call to Write
  // will trigger an implicit WriteHeader(http.StatusOK).
  // Thus explicit calls to WriteHeader are mainly used to
  // send error codes.
  WriteHeader(int)
}

// 实际上传递给 Handler 的对象是:

// A response represents the server side of an HTTP response.
type response struct {
  conn          *conn
  req           *Request // request for this response
  wroteHeader   bool     // reply header has been (logically) written
  wroteContinue bool     // 100 Continue response was written

  w  *bufio.Writer // buffers output in chunks to chunkWriter
  cw chunkWriter
  sw *switchWriter // of the bufio.Writer, for return to putBufioWriter

  // handlerHeader is the Header that Handlers get access to,
  // which may be retained and mutated even after WriteHeader.
  // handlerHeader is copied into cw.header at WriteHeader
  // time, and privately mutated thereafter.
  handlerHeader Header
  ...
  status        int   // status code passed to WriteHeader
    ...
}