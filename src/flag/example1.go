package main

import (
	"flag"
	"fmt"
)

var name = flag.String("name", "World", "A name to say hello to") // 设定-name 参数，如果参数不传值，则使用默认的 World

var spanish bool

func init() {
	flag.BoolVar(&spanish, "spanish", false, "Use Spanish language.") // 设定参数 --spanish=true/false
	flag.BoolVar(&spanish, "s", false, "Use Spanish language.")       // 设定参数 --s=true/false
}

func main() {
	flag.Parse()

	if spanish == true {
		fmt.Printf("Hola %s!\n", *name)
	} else {

		fmt.Printf("Hello %s!\n", *name)
	}

	flag.VisitAll(func(flag *flag.Flag) { // 类似打印你的help 设置显示。也就是将所有前面的设置 以-help的形式打印出来，very nice
		format := "\t-%s: %s (Default: '%s')\n"
		fmt.Printf(format, flag.Name, flag.Usage, flag.DefValue)
	})
}

// ./example1 --spanish=true
// ./example1 --spanish=true -name nice

// Hola nice!
//   -name: A name to say hello to (Default: 'World')
//   -s: Use Spanish language. (Default: 'false')
//   -spanish: Use Spanish language. (Default: 'false')
