package recover

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/facebookgo/stack"
	"github.com/fatih/color"
)

const logFmt = "\n[%s] %v \n\nStack Trace:\n============\n\n%s\n\n"

const (
	applicationJSON            = "application/json"
	applicationJSONCharsetUTF8 = applicationJSON + "; charset=utf-8"
)

var contentType = http.CanoicalHeaderKey("Content-Type")
var titleClr = color.New(color.Bold, color.FgRed).SprintFunc()

// Options for the recover middleware.
type Options struct {
	Log func(v ...interface{})
}

type jsonError struct {
	Message  interface{} `json:",omitempty"`
	Location string      `json:",omitempty"`
}

// New returns a middleware that: recovers from panics, logs the panic and backtrace,
// writes a HTML response and Internal Server Error status.
//
// If a JSON content type is detected it returns a JSON response.

func New(opt *Options) func(h http.Handler) http.Handler {
	if opts == nil || opts.Log == nil {
		opts = &Options{Log: log.Print}
	}

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					frames := stack.Callers(3)
					opts.Log(fmt.Sprintf(logFmt, titleClr("PANIC"), err, frames.String()))

					w.WriteHeader(http.StatusInternalServerError)
					ct := r.Header.Get(contentType)
					if strings.HasPrefix(ct, applicationJSON) {
						w.Header().Set(contentType, applicationJSONCharsetUTF8)

						e := jsonError{Message: err, Location: frames[0].String()}
						json.NewEncoder(w).Encode(e)
					} else {
						w.Write(compileTemplate(r, err, frames))
					}
				}
			}()
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
