// weather-api.go
//
// This is a simple example of how to use the api package.
//
// The getWeather() function takes the NOAA station id for a given location
// (full list at http://w1.weather.gov/xml/current_obs/) and returns the
// current weather conditions as an xml-formatted string.
//
// Inside main(), getWeather is assigned to respond to requests where
// "/weather" is found in the url from the client; it will send its responses
// back in text/xml format, using utf-8, back to the client.
//
// This example server runs on port 9001, and so any request in the form:
//
// http://[localhost/domain/ip of server]:9001/weather?q=[station id]
//
// will work, as long as [station id] corresponds to a valid NOAA value.
//
// Optionally, the request can include an hmac digest for security:
//
// http://[localhost/domain/ip of server]:9001/weather?q=[station id]&d=[hmac
// digest of "q" in sha1 with a shared private key]
//
// The "d" parameter is a sha1 digest of the station id using "secret"
// as the shared private key in this example (in practice, the private
// key is known only by the authorized api client and the server -- see
// http://en.wikipedia.org/wiki/Hmac for more details on how it works).

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dpapathanasiou/go-api"
)

// The getWeather function accepts an http.ResponseWriter and http.Request
// object as input; the latter is used to find specific information about
// the client request, and how to process it.  This function follows
// the GET request pattern, in that it expects data in the url query
// string, even though it does not explicitly check for the GET method
// in the http.Request object.  The http.ResponseWriter is included if
// it's necessary to write additional headers to the reply, beyond the
// Content-type and Content-length values provided automatically by the api
// package (in this specific example, the http.ResponseWriter is not used).
func getWeather(w http.ResponseWriter, r *http.Request) string {
	xml := "<error>Bad Request</error>" // the default response string

	urlValues := r.URL.Query()
	location := urlValues.Get("q")  // required
	apiDigest := urlValues.Get("d") // optional: an hmac digest of "q" (the location) using sha1 with "secret" as the shared private key
	sharedPrivateKey := "secret"    // get this from a database, e.g., IRL

	validRequest := false
	if len(location) == 0 {
		// there is no location ("q") in the request url, so return an error message in xml
		xml = "<error>Please specify a NOAA station id</error>"
	} else {
		// If there is a digest as part of the request, check it using the shared private key
		if len(apiDigest) > 0 {
			if api.DigestMatches(sharedPrivateKey, location, apiDigest) {
				validRequest = true
			} else {
				xml = "<error>Your hmac digest is invalid</error>"
			}
		} else {
			validRequest = true
		}
	}

	if validRequest {
		// Use the location passed as "q" in the request url; it must be a valid NOAA station id,
		// as defined here: http://w1.weather.gov/xml/current_obs/
		res, err := http.Get(fmt.Sprintf("http://w1.weather.gov/xml/current_obs/%s.xml", location))
		if err != nil {
			api.Srv.Logger.Print(err)
		}

		if res.StatusCode == 200 {
			// There is a NOAA station id matching the location in the request url
			// and we were able to get a valid reply from the NOAA server in xml
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				api.Srv.Logger.Print(err)
			}
			res.Body.Close()

			// The xml returned from NOAA has an xsl definition that will confuse some clients, so remove it
			// (probably should use Go's xml parsing package, but this is good enough for a quick-and-dirty example like this
			xml = strings.Replace(string(b), "<?xml-stylesheet href=\"latest_ob.xsl\" type=\"text/xsl\"?>", "", -1)

			// also show an example of what the log output is for this request
			api.Srv.Logger.Printf("INFO: found current weather for NOAA station id %s\n", location)
		} else {
			// While the http.Get() request succeeded, something else went wrong,
			// most likely a 404 status, which means the station id is invalid
			api.Srv.Logger.Printf("WARN: problem finding weather for NOAA station id %s (NOAA server reply: %d)\n", location, res.StatusCode)
			xml = fmt.Sprintf("<error status=\"%d\">Could not get weather for NOAA station id %s</error>", res.StatusCode, location)
		}
	}
	return xml
}

// The main function shows how to use the api package to handle
// different request patterns.  First, a map of type { string:
// func(http.ResponseWriter, *http.Request) } is created.  Next, the map is
// populated with pattern strings (as they as found in the request url),
// mapped to the api.Respond function (which defines both the media type
// and the charset), which calls the function which actually processes
// the client request, and returns a string in the expected format. This
// example defines just one pattern and response (i.e., "/weather" returns
// an xml reply in utf-8), but other patterns and response functions can
// be added to the multiplexer.
func main() {
	handlers := map[string]func(http.ResponseWriter, *http.Request){}
	handlers["/weather"] = func(w http.ResponseWriter, r *http.Request) {
		api.Respond("text/xml", "utf-8", getWeather)(w, r)
	}

	api.NewLocalServer(api.DefaultServerTransport, 9001, api.DefaultServerReadTimeout, false, handlers)
	// To run the api server on a specific IP address, e.g., 192.168.1.1, use NewServer() instead:
	//api.NewServer("192.168.1.1", api.DefaultServerTransport, 9001, api.DefaultServerReadTimeout, false, handlers)

	// Another set of options are the transport layer (default is TCP) and FastCGI

	// To run the api server as a UDP server on a specific IP address or domain, change the transport:
	//api.NewServer("192.168.1.1", "udp", 9001, api.DefaultServerReadTimeout, false, handlers)
	// The same example, but with FastCGI:
	//api.NewServer("192.168.1.1", "udp", 9001, api.DefaultServerReadTimeout, true, handlers)
}
