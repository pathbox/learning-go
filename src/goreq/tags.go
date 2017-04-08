type Place struct {
    Country string `url:"country"`
    City    string `url:"city"`
    ZipCode string `url:"zipcode,omitempty"`
}

type Person struct {
    Place `url:",squash"`

    FirstName string `url:"first_name"`
    LastName  string `url:"last_name"`
    Age       string `url:"age,omitempty"`
    Password  string `url:"-"`
}

johnbull := Person{
  Place: Place{ // squash the embedded struct value
    Country: "UK",
    City:    "London",
    ZipCode: "SW1",
  },
  FirstName: "John",
  LastName:  "Doe",
  Age:       "35",
  Password:  "my-secret", // ignored for parameter
}

goreq.Request{
  Uri:         "http://localhost/",
  QueryString: johnbull,
}.Do()
// =>  `http://localhost/?first_name=John&last_name=Doe&age=35&country=UK&city=London&zip_code=SW1`


// age and zipcode will be ignored because of `omitempty`
// but firstname isn't.
samurai := Person{
  Place: Place{ // squash the embedded struct value
    Country: "Japan",
    City:    "Tokyo",
  },
  LastName: "Yagyu",
}

goreq.Request{
  Uri:         "http://localhost/",
  QueryString: samurai,
}.Do()
// =>  `http://localhost/?first_name=&last_name=yagyu&country=Japan&city=Tokyo`