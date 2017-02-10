package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"
)

type interval []time.Duration

func (i *interval) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *interval) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("interval fla already set")
	}
	for _, dt := range strings.Split(value, ",") {
		duration, err := time.ParseDuration(dt)
		if err != nil {
			return err
		}
		*i = append(*i, duration)
	}
	return nil
}

var intervalFlag interval

func init() {
	flag.Var(&intervalFlag, "deltaT", "comma-separated list of intervals to use between events")
}

func main() {
	flag.Parse()
	fmt.Println(intervalFlag)
}
