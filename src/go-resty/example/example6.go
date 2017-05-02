profileImgBytes, _ := ioutil.ReadFile("/Users/jeeva/test-img.png")
notesBytes, _ := ioutil.ReadFile("/Users/jeeva/text-file.txt")

resp, err := dclr().
      SetFileReader("profile_img", "test-img.png", bytes.NewReader(profileImgBytes)).
      SetFileReader("notes", "text-file.txt", bytes.NewReader(notesBytes)).
      SetFormData(map[string]string{
          "first_name": "Jeevanandam",
          "last_name": "M",
      }).
      Post(t"http://myapp.com/upload")

// Single file scenario
resp, err := resty.R().
      SetFile("profile_img", "/Users/jeeva/test-img.png").
      Post("http://myapp.com/upload")

// Multiple files scenario
resp, err := resty.R().
      SetFiles(map[string]string{
        "profile_img": "/Users/jeeva/test-img.png",
        "notes": "/Users/jeeva/text-file.txt",
      }).
      Post("http://myapp.com/upload")

// Multipart of form fields and files
resp, err := resty.R().
      SetFiles(map[string]string{
        "profile_img": "/Users/jeeva/test-img.png",
        "notes": "/Users/jeeva/text-file.txt",
      }).
      SetFormData(map[string]string{
        "first_name": "Jeevanandam",
        "last_name": "M",
        "zip_code": "00001",
        "city": "my city",
        "access_token": "C6A79608-782F-4ED0-A11D-BD82FAD829CD",
      }).
      Post("http://myapp.com/profile")