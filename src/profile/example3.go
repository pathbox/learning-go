package main

import (
	"fmt"
	"github.com/pkg/profile"
)

func main() {
	// p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	for i := 0; i < 100000; i++ {
		fmt.Println(i)
	}
	defer profile.Start().Stop()
}

//  go build
// $ ./myprogram
// $ go tool pprof —text ./myprogram
// Alternatively you can generate other kinds of files, for example a PDF call-graph:
// $ go tool pprof —pdf ./myprogram /var/folders/jm/j__b36sn04n5__scnnw6bpzr0000gn/T/profile152611925/cpu.pprof > out.pdf
