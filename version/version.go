// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package version has the version information about the application
package version

//Version is the version of application
type Version struct {
	//Code is semivar code of the version
	Code string
	//API is the api version of the application
	API string
}

var (
	//V1 is the version 1 of the application
	V1 = Version{Code: "v1.0.0", API: "v1"}
)

var (
	//Default stores the current version of the application
	Default = V1
)

const (
	//AppName is the name of the application
	AppName = "websockets"
)
