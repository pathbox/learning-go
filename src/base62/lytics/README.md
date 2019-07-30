https://github.com/lytics/base62

```
// Encode a value
urlVal := "http://www.biglongurl.io/?utm_content=content&utm_campaign=campaign"
encodedUrl := base62.StdEncoding.EncodeToString([]byte(urlVal))
	
// Unencoded it
byteUrl, err := base62.StdEncoding.DecodeString(encodedUrl)

```