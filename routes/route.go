// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"context"
	"net/http"

	"github.com/cuttle-ai/websockets/log"
	"github.com/cuttle-ai/websockets/routes/response"

	authConfig "github.com/cuttle-ai/auth-service/config"

	"github.com/cuttle-ai/websockets/version"
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
}

type appCtxKey struct {
	key string
}

//AppContextKey is the key with which the application is saved in the request context
var AppContextKey = appCtxKey{key: "app-context"}

//Register registers the route with the default http handler func
func (r Route) Register(s *http.ServeMux) {
	/*
	 * If the route version is default version then will register it without version string to http handler
	 * Will register the router with the http handler
	 */
	if r.Version == version.Default.API {
		s.Handle(r.Pattern, r)
	}
	s.Handle("/"+r.Version+r.Pattern, r)
}

//ServeHTTP implements HandlerFunc of http package. It makes use of the context of request
func (r Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	/*
	 * Will get the context
	 * We will get the auth-access token from the header
	 * Will get session information about the logged in user
	 * We will fetch the app context for the request
	 * If app contexts have exhausted, we will reject the request
	 * Then we will set the app context in request
	 * Execute request handler func
	 * After execution return the app context
	 */
	//getting the context
	ctx := req.Context()

	//getting the auth token from the header
	cookie, cErr := req.Cookie(authConfig.AuthHeaderKey)
	if cErr != nil {
		log.Warn("Auth cookie not found")
		response.WriteError(res, response.Error{Err: "Couldn't find the auth header " + authConfig.AuthHeaderKey}, http.StatusForbidden)
		_, cancel := context.WithCancel(ctx)
		cancel()
		return
	}

	//will get information about the user
	u, ok := authConfig.GetAutenticatedUser(cookie.Value)
	if !ok {
		log.Warn("User information not found the given auth header")
		response.WriteError(res, response.Error{Err: "Couldn't find the user session " + cookie.Value}, http.StatusForbidden)
		_, cancel := context.WithCancel(ctx)
		cancel()
		return
	}
	sess := authConfig.Session{ID: cookie.Value, Authenticated: true, User: &u}

	//fetching the app context
	appCtxReq := AppContextRequest{
		Type:    Get,
		Out:     make(chan AppContextRequest),
		Session: sess,
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

	//executing the handler
	r.HandlerFunc(c, res, req)

	//cancelling the context
	cancel()
}
