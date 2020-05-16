// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package log is used to print logs based of log types
package log

import (
	"fmt"
	"log"

	"github.com/cuttle-ai/websockets/config"
)

//Log types for logger
const (
	//INFO is for informative logs
	INFO = "INFO"
	//DEBUG is for debugging the app
	DEBUG = "DEBUG"
	//WARN is for warning signatures
	WARN = "WARN"
	//ERROR is for errors
	ERROR = "ERROR"
	//PANIC is for panic log prefix
	PANIC = "PANIC"
)

//Info logs the info logs of the application
func Info(l ...interface{}) {
	log.Print(INFO+": ", fmt.Sprintln(l...))
}

//Debug logs the debug logs of the application if debug logs are not switched off
func Debug(l ...interface{}) {
	//Checking if Debug log is off
	if config.PRODUCTION == 0 {
		return
	}
	log.Print(DEBUG+": ", fmt.Sprintln(l...))
}

//Warn logs the warning logs of the application
func Warn(l ...interface{}) {
	log.Print(WARN+": ", fmt.Sprintln(l...))
}

//Error logs the error logs of the application
func Error(l ...interface{}) {
	log.Print(ERROR+": ", fmt.Sprintln(l...))
}

//Fatal is used to print logs for events which causes the app to exit
func Fatal(l ...interface{}) {
	/*
	 * We will call log.Fatal
	 */
	log.Fatal(PANIC+": ", fmt.Sprintln(l...))
}
