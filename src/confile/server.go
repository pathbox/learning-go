package main

import (
	"log"
	"net/http"

	"github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
)

type globalOptions struct {
	MarginBottom string `json:"margin_bottom,omitempty"` //Set the page bottom margin
	MarginLeft   string `json:"margin_left,omitempty"`   //Set the page left margin (default 10mm)
	MarginRight  string `json:"margin_right,omitempty"`  //Set the page right margin (default 10mm)
	MarginTop    string `json:"margin_top,omitempty"`    //Set the page top margin
	PageHeight   string `json:"page_height,omitempty"`   //Page height
	PageSize     string `json:"page_size,omitempty"`     //Set paper size to: A4, Letter, etc. (default A4)
	PageWidth    string `json:"page_width,omitempty"`    //Page width
	Background   string `json:"background,omitempty"`    //是否带有HTML页面的背景图片或背景色的
	Orientation  string `json:"orientation,omitempty"`   // 横向/纵向 Set orientation to Landscape or Portrait (default Portrait)
}

type Params struct {
	URL    string `json:"url,omitempty"`
	Type   string `json:"type,omitempty"`
	UIID   string `json:"uidID,omitempty"`
	DataID string `json:"dataID,omitempty"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("OK"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/pdf", handleHtmlPDF)
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":9011", router))
}

func handleHtmlPDF(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pa := ParseParams(ps)
	r, bs := pa.ValidateParams()

	if r {
		w.Write(bs)
	}
}

// 解析参数
func ParseParams(ps httprouter.Params) *Params {
	m := make(map[string]string)
	for i := range ps {
		key := ps[i].Key
		m[key] = ps[i].Value
	}

	b, err := jsoniter.Marshal(m)
	if err != nil {
		log.Panic(err)
	}
	pa := &Params{}
	jsoniter.Unmarshal(b, pa)

	return pa
}

// 校验参数 true 表示没错,false 表示有错
func (pa *Params) ValidateParams() (bool, string) {

	if IsBlank(pa.URL) {
		return false, "url must pass"
	}
	if IsBlank(pa.UIID) {
		return false, "uiID must pass"
	}
	if IsBlank(pa.DataID) {
		return false, "dataID must pass"
	}
	if IsBlank(pa.Type) {
		return false, "type must pass"
	}

	return true, ""
}

func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}
	return false
}
