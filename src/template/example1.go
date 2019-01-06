package main

import (
	"os"
	"strings"
	"text/template"
)

func main() {
	termTmpl := template.Must(template.New("termTmpl").Parse(strings.Replace(`
    +---------------------------------------------------------------------+
    |                                                                     |
    |             _o/ Hello {{ .Name }}!
    |                                                                     |
    |                                                                     |
    |  Did you know that ssh sends all your public keys to any server     |
    |  it tries to authenticate to?                                       |
    |                                                                     |
    |  That's how we know you are @{{ .User }} on GitHub!
    |                                                                     |
    |  Ah, maybe what you didn't know is that GitHub publishes all users' |
    |  ssh public keys and Ben (benjojo.co.uk) grabbed them all.          |
    |                                                                     |
    |  That's pretty handy at times :) for example your key is at         |
    |  https://github.com/{{ .User }}.keys
    |                                                                     |
    |                                                                     |
    |  P.S. This whole thingy is Open Source! (And written in Go!)        |
    |  https://github.com/FiloSottile/whosthere                           |
    |                                                                     |
    |  -- @FiloSottile (https://twitter.com/FiloSottile)                  |
    |                                                                     |
    +---------------------------------------------------------------------+
`, "\n", "\n\r", -1)))

	termTmpl.Execute(os.Stdout, struct{ Name, User string }{"Joe", "Cool Guy"})
}
