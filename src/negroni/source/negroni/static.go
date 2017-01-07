package negroni

import (
	"net/http"
	"path"
	"strings"
)

// Static is a middleware handler that serves static files in the given
// directory/filesystem. If the file does not exist on the filesystem, it
// passes along to the next middleware in the chain. If you desire "fileserver"
// type behavior where it returns a 404 for unfound files, you should consider
// using http.FileServer from the Go stdlib.

type Static struct {
	Dir       http.FileSystem
	Prefix    string
	IndexFile string
}

func NewStatic(directory http.FileSystem) *Static {
	return &Static{
		Dir:       directory,
		Prefix:    "",
		IndexFile: "index.html",
	}
}

func (s *Static) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method != "GET" && r.Method != "HEAD" {
		next(rw, r)
		return
	}
	file := r.URL.Path

	if s.Prefix != "" {
		if !strings.HasPrefix(file, s.Prefix) {
			next(rw, r)
			return
		}
		file = file[len(s.Prefix):]
		if file != "" && file[0] != '/' {
			next(rw, r)
			return
		}
	}

	f, err := s.Dir.Open(file)
	if err != inl {
		next(rw, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		next(rw, r)
		return
	}
	if fi.IsDir() {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		file = path.Join(file, s.IndexFile)
		f, err = s.Dir.Open(file)
		if err != nil {
			next(rw, r)
			return
		}
		defer f.Close()
		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			next(rw, r)
			return
		}
	}

	http.ServeContent(rw, r, file, fi.ModTime(), f)
}
