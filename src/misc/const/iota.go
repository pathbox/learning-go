package main

import (
	"fmt"
)

const (
	CategoryBooks    = iota // 0
	CategoryHealth          // 1
	CategoryClothing        // 2
)

type Stereotype int

const (
	TypicalNoob           Stereotype = iota // 0
	TypicalHipster                          // 1
	TypicalUnixWizard                       // 2
	TypicalStartupFounder                   // 3
)

// skipping values
type AudioOutput int

const (
	OutMute   AudioOutput = iota // 0
	OutMono                      // 1
	OutStereo                    // 2
	_
	_
	OutSurround // 5
)

type Allergen int

const (
	IgEggs         Allergen = 1 << iota // 1 << 0 which is 00000001
	IgChocolate                         // 1 << 1 which is 00000010
	IgNuts                              // 1 << 2 which is 00000100
	IgStrawberries                      // 1 << 3 which is 00001000
	IgShellfish                         // 1 << 4 which is 00010000
)

type ByteSize float64

const (
	_           = iota             // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota) // 1 << (10*1)
	MB                             // 1 << (10*2)
	GB                             // 1 << (10*3)
	TB                             // 1 << (10*4)
	PB                             // 1 << (10*5)
	EB                             // 1 << (10*6)
	ZB                             // 1 << (10*7)
	YB                             // 1 << (10*8)
)

func CountAllTheThings(i int) string {
	return fmt.Sprintf("there are %d things", i)
}

func CountAllTheThings1(i Stereotype) string {
	return fmt.Sprintf("there are %d things", i)
}

func main() {
	fmt.Println("CategoryBooks: ", CategoryBooks)
	fmt.Println("CategoryHealth: ", CategoryHealth)
	fmt.Println("CategoryClothing: ", CategoryClothing)

	fmt.Println("TypicalNoob: ", TypicalNoob)
	fmt.Println("TypicalHipster: ", TypicalHipster)
	fmt.Println("TypicalUnixWizard: ", TypicalUnixWizard)
	fmt.Println("TypicalStartupFounder: ", TypicalStartupFounder)

	n := TypicalHipster
	// fmt.Println(CountAllTheThings(n)) // error
	fmt.Println(CountAllTheThings(int(n)))

	fmt.Println(CountAllTheThings1(n))

	//常量在 Go 中是弱类型直到它使用在一个严格的上下文环境中
	fmt.Println(CountAllTheThings(2))
	fmt.Println(CountAllTheThings1(2))

	fmt.Println("OutMute: ", OutMute)
	fmt.Println("OutMono: ", OutMono)
	fmt.Println("OutStereo: ", OutStereo)
	fmt.Println("OutSurround: ", OutSurround)

	fmt.Println("IgEggs: ", IgEggs)
	fmt.Println("IgChocolate: ", IgChocolate)
	fmt.Println("IgNuts: ", IgNuts)
	fmt.Println("IgStrawberries: ", IgStrawberries)
	fmt.Println("IgShellfish: ", IgShellfish)

	fmt.Println("KB: ", KB)
	fmt.Println("MB: ", MB)
	fmt.Println("GB: ", GB)
	fmt.Println("TB: ", TB)
	fmt.Println("PB: ", PB)
	fmt.Println("EB: ", EB)
	fmt.Println("ZB: ", ZB)
	fmt.Println("YB: ", YB)
}
