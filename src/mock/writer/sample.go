package sample

type writer interface {
	Write([]byte) (int, error)
}
