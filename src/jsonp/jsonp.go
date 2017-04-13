// <script src> </script> 这样的标签可以跨域

package jsonp

import (
	"fmt"
	"net/http"
	"net/url"
)

//
// Example of Use
//
// In a handler you build some JSON then call JsonP on the return value.
//

// func handleVersion(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Set("Content-Type", "application/json")
// 	io.WriteString(res, jsonp.JsonP(fmt.Sprintf(`{"status":"success", "version":"1.0.0"}`+"\n"), res, req))
// }

var JSON_Prefix string = ""

func SetJsonPrefix(p string) {
	JSON_Prefix = p
}

// Take a string 's' in JSON and if a get parameter "callback" is specified then format this for JSONP callback.
// If it is not a JSONp call (no "callback" parameter) then add JSON_Prefix to the beginning.

func JsonP(s string, res http.ResponseWriter, req *http.Request) string {
	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return JSON_Prefix + s
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return JSON_Prefix + s
	}
	callback := m.Get("callback")
	if callback != "" {
		res.Header().Set("Content-Type", "application/javascript")
		return fmt.Sprintf("%s(%s);", callback, s)
	} else {
		return JSON_Prefix + s
	}
}

// If "callback" is not "" then convert the JSON string 's' to a JSONp callback.
// If it is not a JSONp call (no "callback" parameter) then add JSON_Prefix to the beginning.
func JsonP_Param(s string, res http.ResponseWriter, callback string) string {
	if callback != "" {
		res.Header().Set("Content-Type", "application/javascript")
		return fmt.Sprintf("%s(%s);", callback, s)
	} else {
		return JSON_Prefix + s
	}
}

func PrependPrefix(s string) string {
	return JSON_Prefix + s
}
