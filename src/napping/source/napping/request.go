package napping

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Params map[string]string

func (p Params) AsUrlValues() url.Values {
	result := url.Values{}
	for key, value := range p {
		result.Set(key, value)
	}
	return result
}

// A Request describes an HTTP request to be executed, data structures into
// which the result will be unmarshaled, and the server's response. By using
// a  single object for both the request and the response we allow easy access
// to Result and Error objects without needing type assertions.

type Request struct {
	Url     string      // Raw URL string
	Method  string      // HTTP method to use
	Params  *url.Values // URL query parameters
	Payload interface{} // Data to JSON-encode and POST

	// Can be set to true if Payload is of type *bytes.Buffer and client wants
	// to send it as-is
	RawPayload bool

	// Result is a pointer to a data structure.  On success (HTTP status < 300),
	// response from server is unmarshaled into Result.
	Result interface{}

	// CaptureResponseBody can be set to capture the response body for external use.
	CaptureResponseBody bool

	// ResponseBody exports the raw response body if CaptureResponseBody is true.
	ResponseBody *bytes.Buffer

	// Error is a pointer to a data structure.  On error (HTTP status >= 300),
	// response from server is unmarshaled into Error.
	Error interface{}

	// Optional
	Userinfo *url.Userinfo
	Header   *http.Header

	// Custom Transport if needed.
	Transport *http.Transport

	// The following fields are populated by Send().
	timestamp time.Time      // Time when HTTP request was sent
	status    int            // HTTP status for executed request
	response  *http.Response // Response object from http package
	body      []byte         // Body of server's response (JSON or otherwise)
}

// A Response is a Request object that has been executed.
type Response Request

func (r *Response) Timestamp() time.Time {
	return r.timestamp
}

func (r *Response) Status() int {
	return r.status
}

func (r *Response) HttpResponse() *http.Response {
	return r.response
}

func (r *Response) Unmarshal(v interface{}) error {
	return json.Unmarshal(r.body, v)
}
