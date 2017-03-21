package httpway


import (
  "io"
  "net/http"
  "sync/atomic"

  "bufio"
  "fmt"
  "github.com/julienschmidt/httprouter"
  "net"
)

var requestId uint64 = 0

// get the context associated with request
func GetContext(r *http.Request) *Context {
  crc, ok := r.Body.(contextReadCloser)
  if !ok {
    return nil
  }

  return crc.ctx()
}

// this is the context that is created for each request
type Context struct {
  data map[string]interface{}
  logger Logger
  session Session
  statusCode int
  transferedBytes uint64
  params *httprouter.Params
  payload interface{}
  handlers *[]Handler
  runNextHandlerIdx int
}

// execute the next middleware

func (c *Context) Next(w http.ResponseWriter, r *http.Request) {
  c.runNextHandlerIdx--

  if c.runNextHandlerIdx < 0 {
    panic("No next middleware, don't call it in final handler"
  }

  (*c.handlers)[c.runNextHandlerIdx](w, r)
}

// set a key on context
func (c *Context) Set(key string, value interface{}) {
  c.data[key] = value
}

// get a key from context and if was set
func (c *Context) GetOk(key string) (value interface{}, found bool) {
  value, found = c.data[key]
  return
}

func (c *Context) Get(key string) interface{} {
  reutrn c.data[key]
}

func (c *Context) Has(key string) bool {
  _, has := c.data[key]
  return has
}

func (c *Context) StatusCode() int {
  return c.statusCode
}

func (c *Context) TransferedBytes() uint64{
  return c.transferedBytes
}

func (c *Context) Log() Logger {
  if c.logger == nil{
    panic("No logger set")
  }
  return c.logger
}

func (c *Context) HasLog() bool {
  if c.logger == nil {
    return false
  }

  return true
}

func (c *Context) Session() Session{
  if c.session == nil {
    panic("No session set")
  }

  return c.session
}

func (c *Context) HasSession() bool {
  if c.session == nil {
    return false
  }
  return true
}

func (c *Context) ParamByName(name string) string {
  return c.params.ByName(name)
}

func (c *Context) Payload() interface{} {
  return c.Payload
}

// create context with middlewares chain for the request
func CreateContext(router *Router, w http.ResponseWriter, r *http.Request, handlersLen *int, pr *httprouter) {
  crc := &contextReadClose{
    ReadCloser: r.Body,
    ctxObj: &Context{
      data: make(map[string]interface{}),
      handlers: handlers,
      runNextHandlerIdx: *handLerslen,
      params: pr,
    },
  }

  crc.ctxObj.logger = &internalLogger{router.Logger, atom.AddUint64(&requestId, 1), ""}

  r.Body = crc
  w = &internalResponseWriter{w, crc.ctxObj}

  if router.SessionManager != nil {
    crc.ctxObj.session = router.SessionManager.Get(w, r, crc.ctxObj.logger)

    if crc.ctxObj.session.Username() != "" {
      crc.ctxObj.logger.(*internalLogger).prefix = crc.ctxObj.session.Username()
    }
  }

  return w
}

type contextReadCloser interface {
  io.ReadCloser
  ctx() *Context
}

type contextReadClose struct {
  io.ReadCloser
  ctxObj *Context
}

func (crc *contextReadClose) ctx() *Context {
  return crc.ctxObj
}

type internalResponseWriter struct {
  http.ResponseWriter
  ctx *Context
}

func (irw *internalResponseWriter) Write(b []byte) (n int, err error) {
  if irw.ctx.statusCode == 0 {
    irw.ctx.statusCode = 200
  }
  n, err = irw.ResponseWriter.Write(b)

  irw.ctx.transferedBytes += uint64(n)
  return
}

func (irw *internalResponseWriter) WriteHeader(status int) {
  if irw.ctx.statusCode == 0 {
    irw.ctx.statusCode = status
  }
  irw.ResponseWriter.WriteHeader(status)
}

func (irw *internalResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
  hij, ok := irw.ResponseWriter.(http.Hijacker)
  if !ok {
    return nil, nil, fmt.Errorf("Unable to Hijack the connection %T doesn't implement Hijack", irw.ResponseWriter)
  }

  return hij.Hijack()
}