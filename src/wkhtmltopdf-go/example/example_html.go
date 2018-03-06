package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/andrewcharlton/wkhtmltopdf-go"
)

const page = `
<html>
  <body>
    <h1>Test Page</h1>

	<p>Path: {{.}}</p>
  </body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.New("page").Parse(page))
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, "Hello World")

	doc := wkhtmltopdf.NewDocument()
	pg, err := wkhtmltopdf.NewPageReader(buf)
	if err != nil {
		log.Fatal("Error reading page buffer")
	}
	doc.AddPages(pg)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="test.pdf"`)
	err = doc.Write(w)
	WkhtmlPDF()
	if err != nil {
		log.Fatal("Error serving pdf")
	}
}

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

func WkhtmlPDF() {

	doc := wkhtmltopdf.NewDocument()
	pg := wkhtmltopdf.NewPage("https://qii404.me/2016/07/22/wkhtmltopdf.html")
	doc.AddPages(pg)

	doc.WriteToFile("html.pdf")
}
