package main

import (
	"net/http"

	"github.com/unrolled/secure" // or "gopkg.in/unrolled/secure.v1"
	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("X-Frame-Options header is now `DENY`."))
	})

	secureMiddleware := secure.New(secure.Options{
		FrameDeny: true,
	})

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(mux)

	n.Run("127.0.0.1:3003")
}

// If you would like to add the above security rules directly to your Nginx configuration, everything is below:

// # Allowed Hosts:
// if ($host !~* ^(example.com|ssl.example.com)$ ) {
//     return 500;
// }

// # SSL Redirect:
// server {
//     listen      80;
//     server_name example.com ssl.example.com;
//     return 301 https://ssl.example.com$request_uri;
// }

// # Headers to be added:
// add_header Strict-Transport-Security "max-age=315360000";
// add_header X-Frame-Options "DENY";
// add_header X-Content-Type-Options "nosniff";
// add_header X-XSS-Protection "1; mode=block";
// add_header Content-Security-Policy "default-src 'self'";
// add_header Public-Key-Pins 'pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"';
