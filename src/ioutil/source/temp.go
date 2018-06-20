import (
	"os"
	"strconv"
	"sync"
	"time"
)

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.GetPid())) // 时间戳+当前进程ID
}

func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]

}

func TempDir(dir, prefix string) (name string, err error) {
	if dir == "" { // 如果 没有指定具体的目录，会在临时目录下生产文件目录
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		try := fliepath.Join(dir, prefix+nextSuffix())
		err = os.Mkdir(try, 0700)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		if os.IsNotExist(err) {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return "", err
			}
		}
		if err == nil {
			name = try
		}
		break
	}
	return
}

func TempFile(dir, prefix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := fliepath.Join(dir, prefix+nextSuffix())
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}