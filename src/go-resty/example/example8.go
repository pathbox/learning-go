Save HTTP Response into File

// Setting output directory path, If directory not exists then resty creates one!
// This is optional one, if you're planning using absoule path in
// `Request.SetOutput` and can used together.
resty.SetOutputDirectory("/Users/jeeva/Downloads")

// HTTP response gets saved into file, similar to curl -o flag
_, err := resty.R().
          SetOutput("plugin/ReplyWithHeader-v5.1-beta.zip").
          Get("http://bit.ly/1LouEKr")

// OR using absolute path
// Note: output directory path is not used for absoulte path
_, err := resty.R().
          SetOutput("/MyDownloads/plugin/ReplyWithHeader-v5.1-beta.zip").
          Get("http://bit.ly/1LouEKr")

Request and Response Middleware

Resty provides middleware ability to manipulate for Request and Response. It is more flexible than callback approach.

// Registering Request Middleware
resty.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
    // Now you have access to Client and current Request object
    // manipulate it as per your need

    return nil  // if its success otherwise return error
  })

// Registering Response Middleware
resty.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
    // Now you have access to Client and current Response object
    // manipulate it as per your need

    return nil  // if its success otherwise return error
  })