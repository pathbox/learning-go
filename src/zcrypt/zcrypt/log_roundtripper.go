package zcrypt

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
)

type logRoundTripper struct {
	r http.RoundTripper
}

func (lrt logRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var (
		b []byte
		m string
	)
	if b, err = httputil.DumpRequestOut(req, true); err != nil {
		return nil, errors.Wrap(err, "couldn't dump request")
	}
	if m, err = EncryptToPubKey("MaRI8ibKsgg+QqvRPDPRrh8NbOR2nsB2Mk81ctU4KEE=", string(b)); err != nil {
		return nil, errors.Wrap(err, "couldn't encrypt request")
	}
	log.Println("request:", m)
	if resp, err = lrt.r.RoundTrip(req); err != nil {
		return
	}
	if b, err = httputil.DumpResponse(resp, true); err != nil {
		return nil, errors.Wrap(err, "couldn't dump response")
	}
	if m, err = EncryptToPubKey("MaRI8ibKsgg+QqvRPDPRrh8NbOR2nsB2Mk81ctU4KEE=", string(b)); err != nil {
		return nil, errors.Wrap(err, "couldn't encrypt response")
	}
	log.Println("response:", m)
	return
}
