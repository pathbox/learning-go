package main

import (
	"fmt"

	"github.com/mediocregopher/lever"
)

func main() {
  f := lever.New("test-app", nil)

  f.Add(lever.Param{Name: "--foo"})
  f.Add(lever.Param{Name: "--flag1", Flag: true})
  f.Add(lever.Param{Name: "--bar", Aliases: []string{"-b"}, Description: "wut"})
  f.Add(lever.Param{Name: "--baz", Aliases: []string{"-c"}, Description: "wut", Default: "wat"})
	f.Add(lever.Param{Name: "--flag2", Flag: true})
	f.Add(lever.Param{
		Name:         "--buz",
		Aliases:      []string{"-d"},
		Description:  "wut",
		DefaultMulti: []string{"a", "b", "c"},
	})
	f.Add(lever.Param{Name: "--byz", DefaultMulti: []string{}})
  // fmt.Println(f)
  // fmt.Println(f.Help())
	fmt.Println(f.Example())
	fmt.Println(f.ParamFlag("--flag1"))
	fmt.Println(f.ParamStrs("--buz"))
}
