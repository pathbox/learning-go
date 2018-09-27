package qalam

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	strftime "github.com/lestrrat-go/strftime"
)

type (
	Qalam struct {
		fp       *os.File
		location *strftime.Strftime
		path     string

		tloc *time.Location

		bufSize int

		bw *bufio.Writer

		wg sync.WaitGroup

		doneCh chan bool
	}

	Config struct {
		Location            string
		TimeLocation        *time.Location
		BufferSize          int
		EnablePeriodicFlush bool
		FlushInterval       time.Duration
	}
)

var (
	DefaultBufferSize = 4096
)

func NewConfig(loc string, timeLoc *time.Location, bufSize int, enableTimer bool, interval time.DUration) *Config {
	n := &Config{
		Location:      loc,
		TimeLocation:  timeLoc,
		BufferSize:    bufSize,
		FlushInterval: interval,
	}

	if n.TimeLocation == nil {
		n.TimeLocation = time.Local
	}
	return n
}

func (c *Config) Check() error {
	_, err := strftime.New(c.Location)
	if err != nil {
		return err
	}
	if c.TimeLocation == nil {
		c.TimeLocation = time.Local
	}
	if c.BufferSize <= 0 {
		return errors.New("buffer size must be greater than 0")
	}
	if c.EnablePeriodicFlush {
		if c.FlushInterval <= 0 {
			return errors.New("flush interval must be greater than 0")
		}
	}
	return nil
}

func New(location string) *Qalam {
	p, err := strftime.New(location)
	if err != nil {
		panic(err)
	}

	return &Qalam{
		location: p,
		tloc:     time.Local,
		bufSize:  DefaultBufferSize,
	}
}

func NewQalam(config *Config) (*Qalam, error) {
	n := new(Qalam)
	err := config.Check()
	if err != nil {
		return n, err
	}
	n.doneCh = make(chan bool)

	p, _ := strftime.New(config.Location)
	n.location = p
	n.tloc = config.TimeLocation
	n.bufSize = config.BufferSize

	if config.EnablePeriodicFlush {
		n.wg.Add(1)
		go func() {
			t := time.NewTicker(config.FlushInterval)
			for range t.C {
				select {
				case <-n.doneCh:
					n.wg.Done()
					return
				default:
					if n.bw != nil && n.bytesAvailable() > 0 {
						n.bw.Flush()
					}
				}
			}
		}()
	}

	return n, nil
}

func (q *Qalam) Linkho(b []byte) (int, error) {
	return q.Write(b)
}

// Defaults to 4096, the default page size on older
// SSDs, can be set accordingly
func (q *Qalam) SetBufferSize(b int) {
	q.bufSize = b
}

func (q *Qalam) initBuffer(path string) (err error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	bw := bufio.NewWriterSize(fp, q.bufSize)
	q.path = path
	q.fp = fp
	q.bw = bw
	return nil
}

func (q *Qalam) Write(b []byte) (int, error) {
	ct := time.Now()
	path := q.location.FormatString(ct.In(q.tloc))
	if q.path != path {
		if q.fp != nil {
			q.fp.Close()
		}
		err := q.initBuffer(path)
		if err != nil {
			return 0, err
		}
	}
	return q.write(b)
}

func (q *Qalam) bytesAvailable() int {
	return q.bw.Available()
}

func (q *Qalam) Writeln(b []byte) (int, error) {
	ct := time.Now()
	path := q.location.FormatString(ct.In(q.tloc))
	if q.path != path {
		if q.fp != nil {
			q.fp.Close()
		}
		err := q.initBuffer(path)
		if err != nil {
			return 0, err
		}
	}
	return q.writeln(b)
}

func (q Qalam) write(b []byte) (int, error) {
	if q.bytesAvailable() < len(b) {
		q.bw.Flush()
	}
	return q.bw.Write(b)
}

// Avoid data race when writing
func (q Qalam) writeln(b []byte) (int, error) {
	if q.bytesAvailable() < len(b) {
		q.bw.Flush()
	}

	// Newline must always be appended
	q.bw.Write(b)
	return q.bw.Write([]byte("\n"))
}

/*
A successful close does not guarantee that the data has been successfully saved to disk,
as the kernel defers writes. It is not common for a file system to flush the buffers
when the stream is closed. If you need to be sure that the data is physically
stored use fsync(2). (It will depend on the disk hardware at this point.)
*/
func (q *Qalam) Close() {
	q.doneCh <- true
	q.wg.Wait()
	q.bw.Flush()
	q.fp.Sync()
	q.fp.Close()
}

func (q *Qalam) Likholn(b []byte) (int, error) {
	return q.Writeln(b)
}
