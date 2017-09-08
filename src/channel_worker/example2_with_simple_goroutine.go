func payloadHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    w.WriteHeader(http.StatusMthodNotAllowed)
    return
  }

  // Read the body into a string for json decoding

  var content = &PayloadCollection{}
  err := json.NewDecoder(io.LimitReader(r.Body, MaxLength)).Decode(&content) // NewDecoder 方法将 bytes 转为json，Decode将json值按照content struct结构，反序列化。相当于将json值传到content struct对象上

  if err != nil {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  for _, payload := range content.Payloads {
    go payload.UploadToS3() // Don't do this  当有成百上千的请求时，goroutine 会爆
  }

  w.WriteHeader(http.StatusOK)
}

