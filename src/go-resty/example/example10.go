// Retries are configured per client
resty.DefaultClient.
    // Set retry count to non zero to enable retries
    SetRetryCount(3).
    // You can override initial retry wait time.
    // Default is 100 milliseconds.
    SetRetryWaitTime(5 * time.Second).
    // MaxWaitTime can be overridden as well.
    // Default is 2 seconds.
    SetRetryMaxWaitTime(20 * time.Second)



resty.DefaultClient.
    AddRetryCondition(
        // Condition function will be provided with *resty.Response as a
        // parameter. It is expected to return (bool, error) pair. Resty will retry
        // in case condition returns true or non nil error.
        func(r *resty.Response) (bool, error) {
            return r.StatusCode() == http.StatusTooManyRequests, nil
        }
    )

// REST mode. This is Default.
resty.SetRESTMode()

// HTTP mode
resty.SetHTTPMode()

// Here you go!
// Client 1
client1 := resty.New()
client1.R().Get("http://httpbin.org")
// ...

// Client 2
client2 := resty.New()
client1.R().Head("http://httpbin.org")
// ...

// Bend it as per your need!!!