// Package fasthttprouter is a trie based high performance HTTP request router.
//
// A trivial example is:
//
// package main

// import (
//     "fmt"
//     "log"
//
//     "github.com/buaazp/fasthttprouter"
//     "github.com/valyala/fasthttp"
// )

// func Index(ctx *fasthttp.RequestCtx) {
//     fmt.Fprint(ctx, "Welcome!\n")
// }

// func Hello(ctx *fasthttp.RequestCtx) {
//     fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
// }

// func main() {
//     router := fasthttprouter.New()
//     router.GET("/", Index)
//     router.GET("/hello/:name", Hello)

//     log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
// }
//
// The router matches incoming requests by the request method and the path.
// If a handle is registered for this path and method, the router delegates the
// request to that function.
// For the methods GET, POST, PUT, PATCH and DELETE shortcut functions exist to
// register handles, for all other methods router.Handle can be used.
//
// The registered path, against which the router matches incoming requests, can
// contain two types of parameters:
//  Syntax    Type
//  :name     named parameter
//  *name     catch-all parameter
//
// Named parameters are dynamic path segments. They match anything until the
// next '/' or the path end:
//  Path: /blog/:category/:post
//
//  Requests:
//   /blog/go/request-routers            match: category="go", post="request-routers"
//   /blog/go/request-routers/           no match, but the router would redirect
//   /blog/go/                           no match
//   /blog/go/request-routers/comments   no match
//
// Catch-all parameters match anything until the path end, including the
// directory index (the '/' before the catch-all). Since they match anything
// until the end, catch-all parameters must always be the final path element.
//  Path: /files/*filepath
//
//  Requests:
//   /files/                             match: filepath="/"
//   /files/LICENSE                      match: filepath="/LICENSE"
//   /files/templates/article.html       match: filepath="/templates/article.html"
//   /files                              no match, but the router would redirect
//
// The value of parameters is inside ctx.UserValue
// To retrieve the value of a parameter:
//  // use the name of the parameter
//  user := ps.UserValue("user")
//

package fasthttprouter

import (
	"github.com/valyala/fasthttp"
	"strings"
)

var (
	defaultContentType = []byte("text/plain; charset=utf-8")
	questionMark       = []byte("?")
)

type Router struct {
	tress map[string]*node

	RedirectTrailingSlash bool

	RedirectFixedPath bool

	HandleMethodNotAllowed bool

	HandleOPTIONS bool

	NotFound fasthttp.RequestHandler

	MethodNotAllowed fasthttp.RequestHandler

	PanicHandler func(*fasthttp.RequestCtx, interface{})
}

func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *Router) GET(path string, handle fasthttp.RequestHandler) {
	r.Handle("GET", path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) HEAD(path string, handle fasthttp.RequestHandler) {
	r.Handle("HEAD", path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) OPTIONS(path string, handle fasthttp.RequestHandler) {
	r.Handle("OPTIONS", path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *Router) POST(path string, handle fasthttp.RequestHandler) {
	r.Handle("POST", path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) PUT(path string, handle fasthttp.RequestHandler) {
	r.Handle("PUT", path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle fasthttp.RequestHandler) {
	r.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle fasthttp.RequestHandler) {
	r.Handle("DELETE", path, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).

func (r *Router) Handle(method, path string, handle fasthttp.RequestHandler) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.tress == nil {
		r.tress = make(map[string]*node)
	}

	root := r.tress[method]
	if root == nil {
		root = new(node)
		r.tress[method] = root
	}

	root.addRoute(path, handle)
}

func (r *Router) ServeFiles(path string, rootPath string) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}
	prefix := path[:len(path)-10]

	fileHandler := fasthttp.FSHandler(rootPath, strings.Count(prefix, "/"))

	r.GET(path, func(ctx *fasthttp.RequestCtx) {
		fileHandler(ctx)
	})
}

func (r *Router) recv(ctx *fasthttp.RequestCtx) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string, ctx *fasthttp.RequestCtx) (fasthttp.RequestHandler, bool) {
	if root := r.tress[method]; root != nil {
		return root.getValue(path, ctx)
	}
	return nil, false
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	if path == "*" || path == "/*" {
		for method := range r.tress {
			if method == "OPTIONS" {
				continue
			}

			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else {
		for method := range r.tress {
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _ := r.tress[method].getValue(path, nil)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

// Handler makes the router implement the fasthttp.ListenAndServe interface
func (r *Router) Handler(ctx *fasthttp.RequestCtx) {
	if r.PanicHandler != nil {
		defer r.recv(ctx)
	}

	path := string(ctx.Path())
	method := string(ctx.Method())
	if root := r.tress[method]; root != nil {
		if f, tsr := root.getValue(path, ctx); f != nil {
			f(ctx)
			return
		} else if method != "CONNECT" && path != "/" {
			code := 301
			if method != "GET" {
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				var uri string
				if len(path) > 1 && path[len(path)-1] == '/' {
					uri = path[:len(path)-1]
				} else {
					uri = path + "/"
				}
				ctx.Redirect(uri, code)
				return
			}
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)

				if found {
					queryBuf := ctx.URI().QueryString()
					if len(queryBuf) > 0 {
						fixedPath = append(fixedPath, questionMark...)
						fixedPath = append(fixedPath, queryBuf...)
					}
					uri := string(fixedPath)
					ctx.Redirect(uri, code)
					return
				}
			}
		}
	}
	if method == "OPTIONS" {
		// Handle OPTIONS requests
		if r.HandleOPTIONS {
			if allow := r.allowed(path, method); len(allow) > 0 {
				ctx.Response.Header.Set("Allow", allow)
				return
			}
		}
	} else {
		// Handle 405
		if r.HandleMethodNotAllowed {
			if allow := r.allowed(path, method); len(allow) > 0 {
				ctx.Response.Header.Set("Allow", allow)
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed(ctx)
				} else {
					ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
					ctx.SetContentTypeBytes(defaultContentType)
					ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
				}
				return
			}
		}
	}

	// Handle 404
	if r.NotFound != nil {
		r.NotFound(ctx)
	} else {
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound),
			fasthttp.StatusNotFound)
	}

}
