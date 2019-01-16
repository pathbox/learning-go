package main

import (
	"os"

	"github.com/mdp/qrterminal"
)

func main() {
	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK, // 这样才能被微信扫描成果
		WhiteChar: qrterminal.WHITE, //
		QuietZone: 1,
	}

	qrterminal.GenerateWithConfig("https://github.com", config)
}
