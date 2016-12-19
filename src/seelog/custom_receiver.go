package main

import (
	seelog "github.com/cihub/seelog"
)

type SomeLogger struct {
	inner seelog.LoggerInterface
}

func (sw *SomeLogger) Debug(s string) {
	sw.inner.Debug(s)

}
func (sw *SomeLogger) Info(s string) {
	sw.inner.Info(s)
}
func (sw *SomeLogger) Error(s string) {
	sw.inner.Error(s)
}

var log = &SomeLogger{}

func init() {
	var err error
	log.inner, err = seelog.LoggerFromConfigAsString(
		`<seelog>
            <outputs>
                <console formatid="fmt"/>
            </outputs>
            <formats>
                <format id="fmt" format="[%Func] [%Lev] %Msg%n"/>
            </formats>
        </seelog>
        `)
	if err != nil {
		panic(err)
	}
	log.inner.SetAdditionalStackDepth(1)
}

func main() {
	defer log.inner.Flush()
	log.Debug("Test")
	log.Info("Test2")
}
