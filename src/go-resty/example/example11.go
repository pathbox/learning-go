restry.SetDebug(true)

// Using your custom log writer
logFile, _ := os.OpenFile("/Users/jeeva/go-resty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

restry.SetLogger(logFile)

restry.SetTLSClientConfig(&tls.Config{RootCAs: roots})


restry.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

restry.SetTimeout(time.Duration(1 * time.Minute))

resty.SetHostURL("http://httpbin.org")

// Headers for all request
resty.SetHeader("Accept", "application/json")
resty.SetHeaders(map[string]string{
        "Content-Type": "application/json",
        "User-Agent": "My custom User Agent String",
      })

// Cookies for all request
resty.SetCookie(&http.Cookie{
      Name:"go-resty",
      Value:"This is cookie value",
      Path: "/",
      Domain: "sample.com",
      MaxAge: 36000,
      HttpOnly: true,
      Secure: false,
    })
resty.SetCookies(cookies)

// URL query parameters for all request
resty.SetQueryParam("user_id", "00001")
resty.SetQueryParams(map[string]string{ // sample of those who use this manner
      "api_key": "api-key-here",
      "api_secert": "api-secert",
    })
resty.R().SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more")

// Form data for all request. Typically used with POST and PUT
resty.SetFormData(map[string]string{
    "access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
  })

// Basic Auth for all request
resty.SetBasicAuth("myuser", "mypass")

// Bearer Auth Token for all request
resty.SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")

// Enabling Content length value for all request
resty.SetContentLength(true)

// Registering global Error object structure for JSON/XML request
resty.SetError(&Error{})    // or resty.SetError(Error{})


// Unix Socket

unixSocket := "unix:///var/run/my_socket.sock"

// Create a Go's http.Transport so we can set it in resty.
transport := http.Transport{
  Dial: func(_, _ string) (net.Conn, error) {
    return net.Dial("unix", unixSocket)
  },
}

// Set the previous transport that we created, set the scheme of the communication to the
// socket and set the unixSocket as the HostURL.
r := resty.New().SetTransport(transport).SetScheme("http").SetHostURL(unixSocket)

// No need to write the host's URL on the request, just the path.
r.R().Get("/index.html")