package bitcask

const (
	defaultMaxFileSize = 1 << 31 // 2G
)

type Option struct {
	Dir         string
	MaxFileSize uint64
	MergeSecs   int
}

func NewOption(dir string, MaxFileSize uint64) Option {
	if MaxFileSize <= 0 {
		MaxFileSize = defaultMaxFileSize
	}
	return Option{
		Dir:         dir,
		MaxFileSize: MaxFileSize,
	}
}
