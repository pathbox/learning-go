package main

import (
	"fmt"

	"os"

	"strings"

	"text/template"
)

type Person struct {
	Name

	string

	Emails []string
}

const templ = `The name is {{.Name}}.

{{range .Emails}}

An email is "{{. | emailExpand}}"

{{end}}

`

func EmailExpander(args ...interface{}) string {

	ok := false

	var s string

	if len(args) == 1 {

		s, ok = args[0].(string)

	}

	if !ok {

		s = fmt.Sprint(args...)

	}

	// find the @ symbol

	substrs := strings.Split(s, "@")

	if len(substrs) != 2 {

		return s

	}

	// replace the @ by " at "

	return (substrs[0] + " at " + substrs[1])

}

func main() {

	person := Person{

		Name: "jan",

		Emails: []string{"jan@newmarch.name", "jan.newmarch@gmail.com"},
	}

	t := template.New("Person template")

	// add our function

	t = t.Funcs(template.FuncMap{"emailExpand": EmailExpander})

	t, err := t.Parse(templ)

	checkError(err)

	err = t.Execute(os.Stdout, person)

	checkError(err)

}

func checkError(err error) {

	if err != nil {

		fmt.Println("Fatal error ", err.Error())

		os.Exit(1)

	}

}
