package main

import (
	"fmt"

	columnize "github.com/ryanuber/columnize"
)

func main() {
	output := []string{
		"Name | Gender | Age",
		"Bob | Male | 38",
		"Sally | Female | 26",
	}
	result := columnize.SimpleFormat(output)

	fmt.Println(result)

	config := columnize.DefaultConfig()
	config.Delim = "|"
	config.Glue = "  "
	config.Prefix = ""
	config.Empty = ""
	config.NoTrim = false

	result = columnize.Format(output, config)
}

/*
config := columnize.DefaultConfig()
config.Delim = "|"
config.Glue = "  "
config.Prefix = ""
config.Empty = ""
config.NoTrim = false

Delim is the string by which columns of input are delimited
Glue is the string by which columns of output are delimited
Prefix is a string by which each line of output is prefixed
Empty is a string used to replace blank values found in output
NoTrim is a boolean used to disable the automatic trimming of input values

*、