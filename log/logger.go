// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

/* This file contains the definitions of logger interface */

//Logger must be implemented by the logger utilities to be an app logger
type Logger struct {
	//ID of the logger
	ID int
}

//NewLogger returns the new logger with ID initiated
func NewLogger(ID int) *Logger {
	return &Logger{ID: ID}
}

//GetID returns the id of the logger
func (lo *Logger) GetID() int {
	return lo.ID
}

//Info logs the informative logs
func (lo *Logger) Info(l ...interface{}) {
	p := append([]interface{}{"ID:", lo.ID}, l...)
	Info(p...)
}

//Debug logs for the debugging logs
func (lo *Logger) Debug(l ...interface{}) {
	p := append([]interface{}{"ID:", lo.ID}, l...)
	Info(p...)
}

//Warn logs the warning logs
func (lo *Logger) Warn(l ...interface{}) {
	p := append([]interface{}{"ID:", lo.ID}, l...)
	Info(p...)
}

//Error logs the error
func (lo *Logger) Error(l ...interface{}) {
	p := append([]interface{}{"ID:", lo.ID}, l...)
	Info(p...)
}

//Fatal logs the fatal issues
func (lo *Logger) Fatal(l ...interface{}) {
	p := append([]interface{}{"ID:", lo.ID}, l...)
	Info(p...)
}
