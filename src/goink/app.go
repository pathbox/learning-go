package GoInk

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
)

type App struct {
	router  *Router
	routerC map[string]*routerCache
	view    *View
	middle  []Handler
	inter   map[string]Handler
	config  *Config
}

// 工厂方法
func New() *App {
	a := new(App)
	a.router = NewRouter()
	a.routerC = make(map[string]*routerCache)
	a.middle = make([]Handler, 0)
	a.inter = make(map[string]Handler)
	a.config, _ = NewConfig("config.json")
	a.view = NewView(a.config.StringOr("app.view_dir", "view"))
	return a
}

// Use adds middleware handlers.
// Middleware handlers invoke before route handler in the order that they are added.

func (app *App) Use(h ...Handler) {
	app.middle = append(app.middle, h...)
}

func (app *App) Config() *Config {
	return app.config
}

func (app *App) View() *View {
	return app.view
}

func (app *App) handler(res http.ResponseWriter, req *http.Request) {
	context := NewContext(app, res, req)
	defer func() {
		e := recover()
		if e == nil {
			context = nil
			return
		}
		context.Body = []byte(fmt.Sprint(e))
		context.Status = 503
		println(string(context.Body))
		debug.PrintStack()
		if _, ok := app.inter["recover"]; ok {
			app.inter["recover"](context)
		}
		if !context.IsEnd {
			context.End()
		}
		context = nil
	}()

	if _, ok := app.inter["static"]; ok {
		app.inter["static"](context)
		if context.IsEnd {
			return
		}
	}
	if len(app.middle) > 0 {
		for _, h := range app.middle {
			h(context)
			if context.IsEnd {
				break
			}
		}
	}
	if context.IsSend {
		return
	}
	var (
		params map[string]string
		fn     []Handler
		url    = req.URL.Path
	)

	if _, ok := app.routerC[url]; ok {
		params = app.routerC[url].param
		fn = app.routerC[url].fn
	} else {
		params, fn = app.router.Find(url, req.Method)
	}

	if params != nil && fn != nil {
		context.routeParams = params
		rc := new(routerCache)
		rc.params = params
		rc.fn = fn
		app.routerC[url] = rc
		for _, r := range fn {
			f(context)
			if context.IsEnd {
				break
			}
		}
		if !context.IsEnd {
			context.End()
		}
	} else {
		println("router is missing at " + req.URL.Path)
		context.Status = 404
		if _, ok := app.inter["notfound"]; ok {
			app.inter["notfound"](context)
			if !context.IsEnd {
				context.End()
			}
		} else {
			context.Throw(404)
		}
	}
	context = nil
}

// ServeHTTP is HTTP server implement method. It makes App compatible to native http handler.
func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.handler(res, req)
}

func (app *App) Run() {
	addr := app.config.StringOr("app.server", "localhost:9090")
	println("http server run at " + addr)
	e := http.ListenAndServe(addr, app)
	panic(e)
}

// Set app config value
func (app *App) Set(key string, v interface{}) {
	app.config.Set("app."+key, v)
}

// Get app config value if only key string given, return string value.
// If fn slice given, register GET handlers to router with pattern string.
func (app *App) Get(key string, fn ...Handler) string {
	if len(fn) > 0 {
		app.router.Get(key, fn...)
		return ""
	}
	return app.config.String("app." + key)
}

// Register POST handlers to router.
func (app *App) Post(key string, fn ...Handler) {
	app.router.Post(key, fn...)
}

// Register PUT handlers to router.
func (app *App) Put(key string, fn ...Handler) {
	app.router.Put(key, fn...)
}

// Register DELETE handlers to router.
func (app *App) Delete(key string, fn ...Handler) {
	app.router.Delete(key, fn...)
}

// Register handlers to router with custom methods and pattern string.
// Support GET,POST,PUT and DELETE methods.
// Usage:
//     app.Route("GET,POST","/test",handler)
//

func (app *App) Route(method string, key string, fn ...Handler) {
	methods := strings.Split(method, ",")
	for _, m := range methods {
		switch m {
		case "GET":
			app.Get(key, fn...)
		case "POST":
			app.Post(key, fn...)
		case "PUT":
			app.Put(key, fn...)
		case "DELETE":
			app.Delete(key, fn...)
		default:
			println("unknow route method " + m)
		}
	}
}

// Register static file handler.
// It's invoked before route handler after middleware handler.
func (app *App) Static(h Handler) {
	app.inter["static"] = h
}

func (app *App) Recover(h Handler) {
	app.inter["recover"] = h
}

func (app *App) NotFound(h Handler) {
	app.inter["notfound"] = h
}
