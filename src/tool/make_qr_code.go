package main

import (
	"os"

	"github.com/mdp/qrterminal"
)

func main() {
	qrterminal.Generate("https://github.com/mdp/qrterminal", qrterminal.L, os.Stdout)
	qrterminal.Generate("https://github.com/mdp/qrterminal", qrterminal.M, os.Stdout)
	qrterminal.Generate("https://github.com/mdp/qrterminal", qrterminal.H, os.Stdout)
}
