package GoInk

import (
	goUrl "net/url"
	"path"
	"regexp"
	"strings"
)

const (
	ROUTER_METHOD_GET    = "GET"
	ROUTER_METHOD_POST   = "POST"
	ROUTER_METHOD_PUT    = "PUT"
	ROUTER_METHOD_DELETE = "DELETE"
)

type Router struct {
	routeSlice []*Route
}

func NewRouter() *Router {
	rt := new(Router)
	rt.routeSlice = make([]*Route, 0)
	return rt
}

func newRoute() *Route {
	route := new(Route)
	route.params = make([]string, 0)
	return route
}

// Get registers GET handlers with pattern string.
func (rt *Router) Get(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_GET
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

// Post registers POST handlers with pattern string.
func (rt *Router) Post(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_POST
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

// Put registers PUT handlers with pattern string.
func (rt *Router) Put(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_PUT
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

// Delete registers DELETE handlers with pattern string.
func (rt *Router) Delete(pattern string, fn ...Handler) {
	route := newRoute()
	route.regex, route.params = rt.parsePattern(pattern)
	route.method = ROUTER_METHOD_DELETE
	route.fn = fn
	rt.routeSlice = append(rt.routeSlice, route)
}

func (rt *Router) parsePattern(pattern string) (regex *regexp.Regexp, params []string) {
	params = make([]string, 0)
	segments := strings.Split(goUrl.QueryEscape(pattern), "%2F")
	for i, v := range segments {
		if strings.HasPrefix(v, "%3A") {
			segments[i] = `([\w-%\]+)`
			params = append(params, strings.TrimPrefix(v, "%3A"))
		}
	}
	regex, _ = regexp.Compile("^" + strings.Join(segments, "/"+"$"))
	return
}

type Route struct {
	regex  *regexp.Regexp
	method string
	params []string
	fn     []Handler
}

type Handler func(context *Context)

type routerCache struct {
	params map[string]string
	fn     []Handler
}

// Find does find matched rule and parse route url, returns route params and matched handlers.
func (rt *Router) Find(url string, method string) (params map[string]string, fn []Handler) {
	sfx := path.Ext(url)
	url = strings.Replace(url, sfx, "", -1)
	// fix path end slash
	url = goUrl.QueryEscape(url)
	if !strings.HasSuffix(url, "%2F") && sfx == "" {
		url += "%2F"
	}
	url = strings.Replace(url, "%2F", "/", -1)
	for _, r := range rt.routeSlice {
		if r.regex.MatchString(url) && r.method == method {
			p := r.regex.FindStringSubmatch(url)
			if len(p) != len(r.params)+1 {
				continue
			}
			params = make(map[string]string)
			for i, n := range r.params {
				params[n] = p[i+1]
			}
			fn = r.fn
			return
		}
	}
	return nil, nil
}
