type PayloadCollection struct {
  WindowsVersion string `json:"version"`
  Token string `json:"token"`
  Payloads []Payload `json:"data"`
}

type Payload struct {
  // [redacted]
}

func (p *Payload) UploadToS3() error {
   storage_path := fmt.Sprintf("%v/%v", p.storageFolder, time.Now().UnixNano())

   bucket := S3Bucket

   b := new(bytess.Buffer)
   encodeErr := json.NewEncoder(b).Encode(payload)
   if encodeErr != nil {
    return encodeErr
   }

   var acl = s3.Private
   var contentType = "application/octet-stream"

   return bucket.PutReader(storage_path, b, int64(b.Len()), contentType, acl, s3.Options{})
}

