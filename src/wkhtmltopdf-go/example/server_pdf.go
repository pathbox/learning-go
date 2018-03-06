package main

import (
	"bytes"
	"log"
	"net/http"

	wkhtmltopdf "github.com/andrewcharlton/wkhtmltopdf-go"
	"github.com/julienschmidt/httprouter"
)

const page = `
<html>
  <body>
	<p>
		<h1>hello</h1>
	</p>
  </body>
</html>`

type Body struct {
	Content string `json:"content"`
}

func handlePDF(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// result, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()

	// if err != nil {
	// 	log.Panic(err)
	// }

	// body := &Body{}
	// err = jsoniter.Unmarshal(result, body)

	// log.Println(body.Content)

	// if len(body.Content) == 0 {
	// 	w.Write([]byte("content can not be blank"))
	// 	return
	// }

	// if err != nil {
	// 	log.Panic(err)
	// }

	buf := &bytes.Buffer{}
	// tmpl := template.Must(template.New("page").Parse(page))

	// tmpl.Execute(buf, "<h1>hello</h1>")

	buf.Write([]byte(page))

	log.Println(buf.String)

	doc := wkhtmltopdf.NewDocument()
	pg, err := wkhtmltopdf.NewPageReader(buf)
	if err != nil {
		log.Fatal("Error reading page buffer")
	}
	doc.AddPages(pg)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="test.pdf"`)
	err = doc.Write(w)
	if err != nil {
		log.Fatal("Error serving pdf")
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("OK"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/pdf", handlePDF)
	log.Println("Server start")
	log.Fatal(http.ListenAndServe(":9011", router))
}
