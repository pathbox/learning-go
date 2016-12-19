// Copyright (c) 2012 - Cloud Instruments Co. Ltd.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"fmt"
	log "github.com/cihub/seelog"
)

func exceptionsMain() {
	defer log.Flush()
	testMinMax()
	testMin()
	testMax()
	testList()
	testFuncException()
	testFileException()
}

func testMinMax() {
	fmt.Println("testMinMax")
	testConfig := `
<seelog type="sync" minlevel="info" maxlevel="error">
  <outputs><console/></outputs>
</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	log.Trace("NOT Printed")
	log.Debug("NOT Printed")
	log.Info("Printed")
	log.Warn("Printed")
	log.Error("Printed")
	log.Critical("NOT Printed")
}

func testMin() {
	fmt.Println("testMin")
	testConfig := `
<seelog type="sync" minlevel="info">
  <outputs><console/></outputs>
</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	log.Trace("NOT Printed")
	log.Debug("NOT Printed")
	log.Info("Printed")
	log.Warn("Printed")
	log.Error("Printed")
	log.Critical("Printed")
}

func testMax() {
	fmt.Println("testMax")
	testConfig := `
<seelog type="sync" maxlevel="error">
  <outputs><console/></outputs>
</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	log.Trace("Printed")
	log.Debug("Printed")
	log.Info("Printed")
	log.Warn("Printed")
	log.Error("Printed")
	log.Critical("NOT Printed")
}

func testList() {
	fmt.Println("testList")
	testConfig := `
<seelog type="sync" levels="info, trace, critical">
  <outputs><console/></outputs>
</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	log.Trace("Printed")
	log.Debug("NOT Printed")
	log.Info("Printed")
	log.Warn("NOT Printed")
	log.Error("NOT Printed")
	log.Critical("Printed")
}

func testFuncException() {
	fmt.Println("testFuncException")
	testConfig := `
<seelog type="sync" minlevel="info">
  <exceptions>
    <exception funcpattern="*main.test*Except*" minlevel="error"/>
  </exceptions>
  <outputs>
    <console/>
  </outputs>
</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	log.Trace("NOT Printed")
	log.Debug("NOT Printed")
	log.Info("NOT Printed")
	log.Warn("NOT Printed")
	log.Error("Printed")
	log.Critical("Printed")

	log.Current.Trace("NOT Printed")
	log.Current.Debug("NOT Printed")
	log.Current.Info("NOT Printed")
	log.Current.Warn("NOT Printed")
	log.Current.Error("Printed")
	log.Current.Critical("Printed")
}

func testFileException() {
	fmt.Println("testFileException")
	testConfig := `
<seelog type="sync" minlevel="info">
  <exceptions>
    <exception filepattern="*main.go" minlevel="error"/>
  </exceptions>
  <outputs>
    <console/>
  </outputs>
</seelog>`

	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	log.Trace("NOT Printed")
	log.Debug("NOT Printed")
	log.Info("NOT Printed")
	log.Warn("NOT Printed")
	log.Error("Printed")
	log.Critical("Printed")
}
