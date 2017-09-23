package main

import (
	"flag"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"github.com/sunfmin/fanout"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var file = flag.String("file", "", "每行一个词语的列表文本文件")
var workers = flag.Int("workers", 60, "并发数")

func main() {
	flag.Parse()

	fmt.Println(*file)
	f, err := os.Open(*file)
	if err != nil {
		return
	}

	fbytes, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	words := strings.Split(string(fbytes), "\n")
	inputs := []interface{}{}

	for _, word := range words {
		inputs = append(inputs, word)
	}

	a := pinying.NewArgs()

	results, err2 := fanout.ParallelRun(*workers, func(input interface{}) (interface{}, error) {
		word := input.(string)
		if strings.TrimSpace(word) == "" {
			return nil, nil
		}

		pya := pinyin.Pinyin(word, a)

		py := ""
		for _, pye := range pya {
			py = py + pye[0]
		}

		pydowncase := strings.ToLower(py)
		domain := pydowncase + ".com"
		outr, err := domainAvailable(word, domain)

		if err != nil {
			fmt.Println("Error: ", err)
			return nil, nil
		}

		if outr.available {
			fmt.Printf("[Ohh Yeah] %s %s\n", outr.word, outr.domain)
		}

		fmt.Printf("\t\t\t %s %s %s\n", outr.word, outr.domain, outr.summary)

		return outr, nil
	}, inputs)

	fmt.Println("Finished ", len(results), ", Error:", err2)
}

type checkResult struct {
	word      string
	domain    string
	available bool
	summary   string
}

func domainAvailable(word string, domain string) (ch checkResult, err error) {
	var summary string
	var output []byte

	ch.word = word
	ch.domain = domain

	cmd := exec.Command("whois", domain)
	output, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	outputstring := string(output)
	if strings.Contains(outputstring, "No match for \"") {
		ch.available = true
		return
	}

	summary = firstLineOf(outputstring, "Registrant Name") + " => "
	summary = summary + firstLineOf(outputstring, "Expiration Date")
	ch.summary = summary
	return
}

func firstLineOf(content string, keyword string) (line string) {
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		if strings.Contains(l, keyword) {
			line = l
			return
		}
	}
	return
}
