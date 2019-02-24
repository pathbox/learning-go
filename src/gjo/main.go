package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Version gjo version info
type Version struct {
	Program     string `json:"program"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Repo        string `json:"repo"`
	Version     string `json:"version"`
}

var (
	array   = flag.Bool("a", false, "creates an array of words")
	pretty  = flag.Bool("p", false, "pretty-prints")
	version = flag.Bool("v", false, "show version")
)

func isRawString(s string) bool {
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		return true
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return true
	}
	return false
}

func parseValue(s string) interface{} {
	if s == "" {
		return nil
	}
	if isRawString(s) {
		return json.RawMessage(s)
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}
	return s
}

func doArray(args []string) (interface{}, error) {
	jsons := []interface{}{}
	for _, value := range args {
		jsons = append(jsons, parseValue(value))
	}
	return jsons, nil
}

func doObject(args []string) (interface{}, error) {
	jsons := make(map[string]interface{}, len(args))
	for _, arg := range args {
		kv := strings.SplitN(arg, "=", 2)
		s := ""
		if len(kv) > 0 {
			s = kv[0]
		}
		if len(kv) != 2 {
			return nil, fmt.Errorf("Argument %q is neither k=v nor k@v", s)
		}
		key, value := kv[0], kv[1]
		jsons[key] = parseValue(value)
	}
	return jsons, nil
}

func doVersion() error {
	enc := json.NewEncoder(os.Stdout)
	if *pretty {
		enc.SetIndent("", "    ")
	}
	return enc.Encode(&Version{
		Program:     "gjo",
		Description: "This is inspired by jpmens/jo",
		Author:      "gorilla0513",
		Repo:        "https://github.com/skanehira/gjo",
		Version:     "1.0.2",
	})
}

func run() int {
	flag.Parse()
	if *version {
		err := doVersion()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return 2
	}

	var value interface{}
	var err error

	if *array {
		value, err = doArray(args)
	} else {
		value, err = doObject(args)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	enc := json.NewEncoder(os.Stdout)
	if *pretty {
		enc.SetIndent("", "    ")
	}
	err = enc.Encode(value)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
