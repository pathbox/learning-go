package GoInk

import (
	"bytes"
	"html/template"
	"os"
	"path"
	"strings"
)

type View struct {
	Dir           string
	FuncMap       template.FuncMap
	IsCache       bool
	templateCache map[string]*template.Template
}

func (v *View) getTemplateInstance(tpl []string) (*template.Template, error) {
	key := strings.Join(tpl, "-")

	if v.IsCache {
		if v.templateCache[key] != nil {
			return v.templarteCache[key], nil
		}
	}

	var (
		t    *template.Template
		e    error
		file []string = make([]string, len(tpl))
	)
	for i, tp := range tpl {
		file[i] = path.Join(v.Dir, tp)
	}

	t = template.New(path.Base(tpl[0]))
	t.Funcs(v.FuncMap)
	t, e = t.ParseFiles(file...)
	if e != nil {
		return nil, e
	}
	if v.IsCache {
		v.templateCache[key] = t
	}
	return t, nil

}

func (v *View) Render(tpl string, data map[string]interface{}) ([]byte, error) {
	t, e := v.getTemplateInstance(strings.Split(tpl, ","))
	if e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	e = t.Execute(&buf, data)
	if e != nil {
		return nil, e
	}
	return buf.Bytes(), nil
}

func (v *View) Has(tpl string) bool {
	f := path.Join(v.Dir, tpl)
	_, e := os.Stat(f)
	return e == nil
}

func (v *View) NoCache() {
	v.IsCache = false
	v.templateCache = make(map[string]*template.Template)
}

func NewView(dir string) *View {
	v := new(View)
	v.Dir = dir
	v.FuncMap = make(template.FuncMap)
	v.FuncMap["Html"] = func(str string) template.HTML {
		return template.HTML(str)
	}
	v.IsCache = false
	v.templateCache = make(map[string]*template.Template)
	return v
}
