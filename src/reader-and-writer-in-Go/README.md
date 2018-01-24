https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185


type Reader interface {
  Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}