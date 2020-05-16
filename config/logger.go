// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

/* This file contains the definitions of logger interface */

//Logger must be implemented by the logger utilities to be an app logger
type Logger interface {
	//Info logs the informative logs
	Info(l ...interface{})
	//Debug logs for the debugging logs
	Debug(l ...interface{})
	//Warn logs the warning logs
	Warn(l ...interface{})
	//Error logs the error
	Error(l ...interface{})
	//Fatal logs the fatal issues
	Fatal(l ...interface{})
	//GetID returns the ID of the logger
	GetID() int
}
