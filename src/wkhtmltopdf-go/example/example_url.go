package main

import (
	"net/http"

	"github.com/andrewcharlton/wkhtmltopdf-go"
)

func handler(w http.ResponseWriter, r *http.Request) {

	WkhtmlPDF()
}

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

func WkhtmlPDF() {

	doc := wkhtmltopdf.NewDocument()
	pg := wkhtmltopdf.NewPage("https://qii404.me/2016/07/22/wkhtmltopdf.html")
	doc.AddPages(pg)

	doc.WriteToFile("/home/user/html.pdf")
}
