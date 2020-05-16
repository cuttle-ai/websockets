// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package routes has the routes supported by the api with proper versioning done
//Suppose a route is /list, it belonged to v2 and current version is v2. Then route will be available as
// /list and /v2/list. If the current version is not v2 then the api will be exposed only as /list. For using routes
//with a server invoke the InitRoutes function.
package routes

import "net/http"

//routes has the list of routes in the application
var routes = []Route{}

//AddRoutes adds the routes to the routes variable
func AddRoutes(r ...Route) {
	routes = append(routes, r...)
}

//InitRoutes initializes the routes in the application
func InitRoutes(s *http.ServeMux) {
	/*
	 * Will register the routes
	 */
	for _, v := range routes {
		v.Register(s)
	}
}
