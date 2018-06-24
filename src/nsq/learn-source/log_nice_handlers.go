import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func V1(f APIHandler) APIHandler {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
		data, err := f(w, req, ps)
		if err != nil {
			RespondV1(w, err.(Err).Code, err)
			return nil, nil
		}
		RespondV1(w, 200, data)
		return nil, nil
	}
}

func RespondV1(w http.ResponseWriter, code int, data interface{}) {
	var response []byte
	var err error
	var isJSON bool

	if code == 200 {
		switch data.(type) {
		case string:
			response = []byte(data.(string))
		case []byte:
			response = data.([]byte)
		case nil:
			response = []byte{}
		default:
			isJSON = true
			response, err = json.Marshal(data)
			if err != nil {
				code = 500
				data = err
			}
		}
	}

	if code != 200 {
		isJSON = true
		response = []byte(fmt.Sprintf(`{"message":"%s"}`, data))
	}

	if isJSON {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.Header().Set("X-NSQ-Content-Type", "nsq; version=1.0")
	w.WriteHeader(code)
	w.Write(response)
}

type Decorator func(APIHandler) APIHandler

type APIHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)

func Log(logf lg.AppLogFunc) Decorator {
	return func(f APIHandler) APIHandler {
		return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			start := time.Now()
			response, err := f(w, req, ps) // 执行f APIHander，真正的执行逻辑位置
			elapsed := time.Since(start)   // 得到f 执行的耗时
			status := 200
			if e, ok := err.(Err); ok {
				status = e.Code
			}
			logf(lg.INFO, "%d %s %s (%s) %s",
				status, req.Method, req.URL.RequestURI(), req.RemoteAddr, elapsed) // 打印日志信息
			return response, err
		}
	}
}

func LogPanicHandler(logf lg.AppLogFunc) func(w http.ResponseWriter, req *http.Request, p interface{}) {
	return func(w http.ResponseWriter, req *http.Request, p interface{}) {
		logf(lg.ERROR, "panic in HTTP handler - %s", p)
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{500, "INTERNAL_ERROR"}
		}, Log(logf), V1)(w, req, nil)
	}
}

func LogNotFoundHandler(logf lg.AppLogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{404, "NOT_FOUND"}
		}, Log(logf), V1)(w, req, nil)
	})
}

func LogMethodNotAllowedHandler(logf lg.AppLogFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{405, "METHOD_NOT_ALLOWED"}
		}, Log(logf), V1)(w, req, nil)
	})
}
