package main

import (
	"flag"
	"io"
	"log"
	"os"

	cov "github.com/matm/gocov-html/cov"
)

func main() {
	var r io.Reader
	log.SetFlags(0)

	var s = flag.String("s", "", "path to custom CSS file")
	flag.Parse()

	switch flag.NArg() {
	case 0:
		r = os.Stdin
	case 1:
		var err error
		if r, err = os.Open(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Usage: %s data.json\n", os.Args[0])
	}

	if err := cov.HTMLReportCoverage(r, *s); err != nil {
		log.Fatal(err)
	}
}
