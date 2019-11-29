package main

import (
	"flag"
	"fmt"
	"os"

	redis "github.com/pascaldekloe/redis"
)

var (
	addrFlag = flag.String("addr", "localhost:6379", "Redis service address.")
	rawFlag  = flag.Bool("raw", false, "Prints the value as is, instead of a quoted string plus line feed.")
)

var Redis *redis.Client // 全局的redis client
func main() {
	flag.Parse()

	Redis = redis.NewClient(*addFlag, 0, 0)
	defer Redis.Close()

	print(flag.Args())
}

func print(keys []string) {
	for _, key := range keys {
		value, err := Redis.GET(key)
		switch {
		case err != nil:
			os.Stderr.WriteString(err.Error())
			os.Exit(255)
		case *rawFlag:
			os.Stdout.Write(value)
		default:
			fmt.Printf("%q\n", value)
		}
	}
}
