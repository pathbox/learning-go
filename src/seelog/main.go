package main

// go get -u github.com/cihub/seelog
import (
	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	log.Info("Hello from Seelog!")
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		return err
	}
	log.ReplaceLogger(logger)
}

package main

import (
    log "github.com/cihub/seelog"
    "time"
    "errors"
)

type inputData struct {
    x, y int
}

type outputData struct {
    result int
    error bool
}

var inputs chan inputData
var outputs chan outputData
var criticalChan chan int

func internalCalculationFunc(x, y int) (result int, err error) {
    log.Debugf("calculating z. x:%d y:%d", x, y)
    z := y
    switch {
    case x == 3 :
        log.Trace("x == 3")
        panic("Failure.")
    case y == 1 :
        log.Trace("y == 1")
        return 0, errors.New("Error!")
    case y == 2 :
        log.Trace("y == 2")
        z = x
    default :
        log.Trace("default")
        z += x
    }
    log.Tracef("z:%d",z)
    retVal := z-3
    log.Debugf("Returning %d", retVal)

    return retVal, nil
}

func generateInputs(dest chan inputData) {
    time.Sleep(1e9)
    log.Debug("Sending 2 3")
    dest <- inputData{x : 2, y : 3}

    time.Sleep(1e9)
    log.Debug("Sending 2 1")
    dest <- inputData{x : 2, y : 1}

    time.Sleep(1e9)
    log.Debug("Sending 3 4")
    dest <- inputData{x : 3, y : 4}

    time.Sleep(1e9)
    log.Debug("Sending critical")
    criticalChan <- 1
}

func consumeResults(res chan outputData) {
    for {
        select {
            case <- outputs:
            // At this point we get and consume resultant value
        }
    }
}

func processInput(input inputData) {
    defer func() {
        if r := recover(); r != nil {
            log.Errorf("Unexpected error occurred: %v", r)
            outputs <- outputData{result : 0, error : true}
        }
    }()
    log.Infof("Received input signal. x:%d y:%d", input.x, input.y)

    res, err := internalCalculationFunc(input.x, input.y)
    if err != nil {
        log.Warnf("Error in calculation: %s", err.Error())
    }

    log.Infof("Returning result: %d error: %t", res, err != nil)
    outputs <- outputData{result : res, error : err != nil}
}

func main() {
    inputs = make(chan inputData)
    outputs = make(chan outputData)
    criticalChan = make(chan int)
    log.Info("App started.")

    go consumeResults(outputs)
    log.Info("Started receiving results.")

    go generateInputs(inputs)
    log.Info("Started sending signals.")

    for {
        select {
            case input := <- inputs:
                processInput(input)
            case <- criticalChan:
                log.Critical("Caught value from criticalChan: Go shut down.")
                panic("Shut down due to critical fault.")
        }
    }
}