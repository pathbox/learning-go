package main

import (
	"log"
	"os"
)

func main() {
  fileName := "debug.log"
  logFile, err := os.Create(fileName)
  defer logFile.Close()

  if err != nil {
    log.Fatalln("open file error")
  }

  debugLog := log.New(logFile, "[Debug]", log.Llongfile)
  debugLog.Println("A debug message here")
  debugLog.SetPrefix("[Info]")
  debugLog.Println("An Info Message here")
	debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
	debugLog.Println("A different prefix")
}
