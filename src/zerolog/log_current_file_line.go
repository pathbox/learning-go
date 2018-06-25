package main

import (
	"path"
	"runtime"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

func main() {

	Log()
}

func Log() {
	service := "myservice"
	pc, fileName, lineNum, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	filePath := path.Base(fileName)

	zerolog.TimeFieldFormat = ""
	log.Error().
		Str("error", "err msg").
		Msg("error message")
	log.Info().
		Str("service", service).
		Str("file_name", filePath).
		Int("line", lineNum).
		Str("func_name", funcName).
		Msgf("Cannot start %s", service)
}

/*

runtime.Caller(1) 还是 runtime.Caller(2) 应该是看runtime.Caller()所在方法的深度

runtime.Caller()最好是放一个方法中，直接放main中，并不能得到想要的信息

*/
