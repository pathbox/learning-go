// just mentioning about POST as an example with simple flow
// User Login
resp, err := resty.R().
      SetFormData(map[string]string{
        "username": "jeeva",
        "password": "mypass",
      }).
      Post("http://myapp.com/login")

// Followed by profile update
resp, err := resty.R().
      SetFormData(map[string]string{
        "first_name": "Jeevanandam",
        "last_name": "M",
        "zip_code": "00001",
        "city": "new city update",
      }).
      Post("http://myapp.com/profile")

// Multi value form data
criteria := url.Values{
  "search_criteria": []string{"book", "glass", "pencil"},
}
resp, err := resty.R().
      SetMultiValueFormData(criteria).
      Post("http://myapp.com/search")