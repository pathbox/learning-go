package grequests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Response is what is returned to a user when they fire off a request
type Response struct {
	Ok          bool
	Error       error
	RawResponse *http.Response

	StatusCode int

	Header http.Header

	internalByteBuffer *bytes.Buffer
}

func buildResponse(resp *http.Response, err error) (*Response, error) {
	if err != nil {
		return &Response{Error: err}, err
	}

	goodResp := &Response{
		Ok:                 resp.StatusCode >= 200 && resp.StatusCode < 300,
		Error:              nil,
		RawResponse:        resp,
		StatusCode:         resp.StatusCode,
		Header:             resp.Header,
		internalByteBuffer: bytes.NewBuffer([]byte{}),
	}

	return goodResp, nil
}

func (r *Response) Read(p []byte) (n int, err error) {
	if r.Error != nil {
		return -1, r.Error
	}

	return r.RawResponse.Body.Read(p)
}

func (r *Response) Close() error {
	if r.Error != nil {
		return r.Error
	}
	io.Copy(ioutil.Discard, r)

	return r.RawResponse.Body.Close()
}

func (r *Response) DownloadToFile(fileName string) error {
	if r.Error != nil {
		return r.Error
	}

	fd, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer r.Close()
	defer fd.Close()

	if _, err := io.Copy(fd, r.getInternalReader()); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// getInternalReader because we implement io.ReadCloser and optionally hold a large buffer of the response (created by
// the user's request)
func (r *Response) getInternalReader() io.Reader {

	if r.internalByteBuffer.Len() != 0 {
		return r.internalByteBuffer
	}
	return r
}

// XML is a method that will populate a struct that is provided `userStruct` with the XML returned within the
// response body
func (r *Response) XML(userStruct interface{}, charsetReader XMLCharDecoder) error {

	if r.Error != nil {
		return r.Error
	}

	xmlDecoder := xml.NewDecoder(r.getInternalReader())

	if charsetReader != nil {
		xmlDecoder.CharsetReader = charsetReader
	}

	defer r.Close()

	if err := xmlDecoder.Decode(&userStruct); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// JSON is a method that will populate a struct that is provided `userStruct` with the JSON returned within the
// response body
func (r *Response) JSON(userStruct interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	jsonDecoder := json.NewDecoder(r.getInternalReader())
	defer r.Close()

	if err := jsonDecoder.Decode(&userStruct); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// createResponseBytesBuffer is a utility method that will populate the internal byte reader – this is largely used for .String()
// and .Bytes()
func (r *Response) populateResponseByteBuffer() {

	// Have I done this already?
	if r.internalByteBuffer.Len() != 0 {
		return
	}

	defer r.Close()

	// Is there any content?
	if r.RawResponse.ContentLength == 0 {
		return
	}

	// Did the server tell us how big the response is going to be?
	if r.RawResponse.ContentLength > 0 {
		r.internalByteBuffer.Grow(int(r.RawResponse.ContentLength))
	}

	if _, err := io.Copy(r.internalByteBuffer, r); err != nil && err != io.EOF {
		r.Error = err
		r.RawResponse.Body.Close()
	}

}

// Bytes returns the response as a byte array
func (r *Response) Bytes() []byte {

	if r.Error != nil {
		return nil
	}

	r.populateResponseByteBuffer()

	// Are we still empty?
	if r.internalByteBuffer.Len() == 0 {
		return nil
	}
	return r.internalByteBuffer.Bytes()

}

// String returns the response as a string
func (r *Response) String() string {
	if r.Error != nil {
		return ""
	}

	r.populateResponseByteBuffer()

	return r.internalByteBuffer.String()
}

// ClearInternalBuffer is a function that will clear the internal buffer that we use to hold the .String() and .Bytes()
// data. Once you have used these functions – you may want to free up the memory.
func (r *Response) ClearInternalBuffer() {

	if r == nil || r.internalByteBuffer == nil {
		return
	}

	r.internalByteBuffer.Reset()
}
