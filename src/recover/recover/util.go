package recover

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/facebookgo/stack"
)

var gopaths []string

func init() {
	for _, p := range strings.Split(os.Getenv("GOPATH"), ":") {
		if p != "" {
			gopaths = append(gopaths, filepath.Join(p, "src")+"/")
		}
	}
}

func mustReadLines(filename string) []string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return strings.Split(string(bytes), "\n")
}

func appendGOPATH(file string) string {
	for _, p := range gopaths {
		f := filepath.Join(p, file)
		if _, err := os.Stat(f); err == nil {
			return f
		}
	}
	return file
}

func compileTemplate(r *http.Request, err interface{}, frames []stack.Frame) []byte {
	file := appendGOPATH(frames[0].File)
	src := mustReadLines(file)
	start := (frames[0].Line - 1) - 5
	end := frames[0].Line + 5
	lines := src[start:end]

	var buf bytes.Buffer
	t := template.Must(template.New("recover").Parse(panicHTML))
	t.Execute(&buf, struct {
		URL         string
		Err         interface{}
		Name        string
		File        string
		StartLine   int
		SourceLines string
		Frames      []stack.Frame
	}{r.URL.Path, err, frames[0].Name, frames[0].File, start + 1, strings.Join(lines, "\n"), frames})
	return bug.Bytes()
}
