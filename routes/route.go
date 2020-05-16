// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"context"
	"net/http"

	"github.com/cuttle-ai/websockets/log"
	"github.com/cuttle-ai/websockets/routes/response"

	"github.com/cuttle-ai/websockets/version"

	"github.com/cuttle-ai/websockets/config"
)

/*
 * This file has the definition of route data structure
 */

//HandlerFunc is the Handler func with the context
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

//Route is a route with explicit versions
type Route struct {
	//Version is the version of the route
	Version string
	//Pattern is the url pattern of the route
	Pattern string
	//HandlerFunc is the handler func of the route
	HandlerFunc HandlerFunc
	//ParseForm will do a form parse before invoking the handler
	ParseForm bool
}

//AppContextKey is the key with which the application is saved in the request context
const AppContextKey = "app-context"

//Register registers the route with the default http handler func
func (r Route) Register(s *http.ServeMux) {
	/*
	 * If the route version is default version then will register it without version string to http handler
	 * Will register the router with the http handler
	 */
	if r.Version == version.Default.API {
		s.Handle(r.Pattern, http.TimeoutHandler(r, config.ResponseTimeout, "timeout"))
	}
	s.Handle("/"+r.Version+r.Pattern, http.TimeoutHandler(r, config.ResponseTimeout, "timeout"))
}

//ServeHTTP implements HandlerFunc of http package. It makes use of the context of request
func (r Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	/*
	 * Will get the context
	 * Will parse the form
	 * We will fetch the app context for the request
	 * If app contexts have exhausted, we will reject the request
	 * Then we will set the app context in request
	 * Execute request handler func
	 * After execution return the app context
	 */
	//getting the context
	ctx := req.Context()

	//parsing the form
	if r.ParseForm {
		err := req.ParseForm()
		if err != nil {
			//error while parsing the form
			log.Error("Error while parsing the request form", err)
			response.WriteError(res, response.Error{Err: "Couldn't parse the request form"}, http.StatusUnprocessableEntity)
			_, cancel := context.WithCancel(ctx)
			cancel()
			return
		}
	}

	//fetching the app context
	appCtxReq := AppContextRequest{
		Type: Get,
		Out:  make(chan AppContextRequest),
	}
	go SendRequest(AppContextRequestChan, appCtxReq)
	resCtx := <-appCtxReq.Out

	//checking whether the app context exhausted or not
	if resCtx.Exhausted {
		//reject the request
		log.Error("We have exhausted the request limits")
		response.WriteError(res, response.Error{Err: "We have exhuasted the server request limits. Please try after some time."}, http.StatusTooManyRequests)
		_, cancel := context.WithCancel(ctx)
		cancel()
		return
	}

	//setting the app context
	newCtx := context.WithValue(ctx, AppContextKey, resCtx.AppContext)

	//executing the request
	r.Exec(newCtx, res, req)

	//returning the app context
	appCtxReq = AppContextRequest{
		Type:       Finished,
		AppContext: resCtx.AppContext,
	}
	go SendRequest(AppContextRequestChan, appCtxReq)
}

//Exec will execute the handler func. By default it will set response content type as as json.
//It will also cancel the context at the end. So no need of explicitly invoking the same in the handler funcs
func (r Route) Exec(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	/*
	 * Will get the cancel for the context
	 * Will set the content type of response as json
	 * Will execute the handlerfunc
	 * Cancelling the context at the end
	 */
	//getting the context cancel
	c, cancel := context.WithCancel(ctx)

	//setting the content type as json
	res.Header().Set("Content-Type", "application/json")

	//executing the handler
	r.HandlerFunc(c, res, req)

	//cancelling the context
	cancel()
}
