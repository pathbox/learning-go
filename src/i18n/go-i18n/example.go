package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var page = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<body>
<h1>{{ .Title }}</h1>
{{range .Paragraphs}}<p>{{ . }}</p>{{end}}
</body>
</html>
`))

func main() {

	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, lang := range []string{"en", "el"} {
		bundle.MustLoadMessageFile(fmt.Sprintf("active.%v.json", lang))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		name := r.FormValue("name")
		if name == "" {
			name = "Alex"
		}

		myInvoicesCount := 10

		helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "HelloPerson",
			},
			TemplateData: map[string]interface{}{
				"Name": name,
			},
		})

		myInvoices := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: "invoices",
			},
			TemplateData: map[string]interface{}{
				"Count": myInvoicesCount,
			},
			PluralCount: myInvoicesCount,
		})

		err := page.Execute(w, map[string]interface{}{
			"Title": helloPerson,
			"Paragraphs": []string{
				myInvoices,
			},
		})
		if err != nil {
			panic(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
