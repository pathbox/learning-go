type Decorator func(APIHandler) APIHandler

type APIHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)


Decorate(f APIHandler, ds ...Decorator) httprouter.Handle{
	decorated := f
	for _, decorate := ds {
		decorated = decorate(decorated)
	}
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		decorated(w, req, ps)
	}
}

/*
Decorator 接收 APIHandler，不断的将decorated 传入 decorate，有点类似中间件的不断包裹 decorated 和decorate是两种不同类型，decorated作为decorate的参数传入，又返回了decorated类型，所以可以构成循环的包裹wrapper功能
*/