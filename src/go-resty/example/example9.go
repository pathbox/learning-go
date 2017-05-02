// Assign Client Redirect Policy. Create one as per you need
resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))

// Wanna multiple policies such as redirect count, domain name check, etc
resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20),
                        resty.DomainCheckRedirectPolicy("host1.com", "host2.org", "host3.net"))

// Using raw func into resty.SetRedirectPolicy
resty.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
  // Implement your logic here

  // return nil for continue redirect otherwise return error to stop/prevent redirect
  return nil
}))

//---------------------------------------------------

// Using struct create more flexible redirect policy
type CustomRedirectPolicy struct {
  // variables goes here
}

func (c *CustomRedirectPolicy) Apply(req *http.Request, via []*http.Request) error {
  // Implement your logic here

  // return nil for continue redirect otherwise return error to stop/prevent redirect
  return nil
}

// Registering in resty
resty.SetRedirectPolicy(CustomRedirectPolicy{/* initialize variables */})


// Custom Root certificates, just supply .pem file.
// you can add one or more root certificates, its get appended
resty.SetRootCertificate("/path/to/root/pemFile1.pem")
resty.SetRootCertificate("/path/to/root/pemFile2.pem")
// ... and so on!

// Adding Client Certificates, you add one or more certificates
// Sample for creating certificate object
// Parsing public/private key pair from a pair of files. The files must contain PEM encoded data.
cert1, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
if err != nil {
  log.Fatalf("ERROR client certificate: %s", err)
}
// ...

// You add one or more certificates
resty.SetCertificates(cert1, cert2, cert3)

// Setting a Proxy URL and Port
resty.SetProxy("http://proxyserver:8888")

// Want to remove proxy setting
resty.RemoveProxy()

// Set proxy for current request
resp, err := c.R().
    SetProxy("http://sampleproxy:8888").
    Get("http://httpbin.org/get")