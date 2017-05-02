// Note: This is one sample of PUT method usage, refer POST for more combination

// Request goes as JSON content type
// No need to set auth token, error, if you have client level settings
resp, err := resty.R().
      SetBody(Article{
        Title: "go-resty",
        Content: "This is my article content, oh ya!",
        Author: "Jeevanandam M",
        Tags: []string{"article", "sample", "resty"},
      }).
      SetAuthToken("C6A79608-782F-4ED0-A11D-BD82FAD829CD").
      SetError(&Error{}).       // or SetError(Error{}).
      Put("https://myapp.com/article/1234")

// Request goes as JSON content type
// No need to set auth token, error, if you have client level settings
resp, err := resty.R().
      SetBody(Article{
        Tags: []string{"new tag1", "new tag2"},
      }).
      SetAuthToken("C6A79608-782F-4ED0-A11D-BD82FAD829CD").
      SetError(&Error{}).       // or SetError(Error{}).
      Patch("https://myapp.com/articles/1234")