package napping

import (
	"net/url"
)

func Send(r *Request) (*Response, error) {
	s := Sesion{}
	return s.Send(r)
}

func Get(url string, p *url.Values, result, errMsg interface{}) (*Response, error) {
	s := Session{}
	return s.Get(url, p, result, errMsg)
}

func Options(url string, result, errMsg interface{}) (*Response, error) {
	s := Session{}
	return s.Options(url, result, errMsg)
}

func Head(url string, result, errMsg interface{}) (*Response, error) {
	s := Session{}
	return s.Head(url, result, errMsg)
}

// Post sends a POST request.
func Post(url string, payload, result, errMsg interface{}) (*Response, error) {
	s := Session{}
	return s.Post(url, payload, result, errMsg)
}

// Put sends a PUT request.
func Put(url string, payload, result, errMsg interface{}) (*Response, error) {
	s := Session{}
	return s.Put(url, payload, result, errMsg)
}

// Patch sends a PATCH request.
func Patch(url string, payload, result, errMsg interface{}) (*Response, error) {
	s := Session{}
	return s.Patch(url, payload, result, errMsg)
}

// Delete sends a DELETE request.
func Delete(url string, p *url.Values, result, errMsg interface{}) (*Response, error) {
	s := Session{}
	return s.Delete(url, p, result, errMsg)
}
