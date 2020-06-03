// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"time"

	authConfig "github.com/cuttle-ai/auth-service/config"
	"github.com/cuttle-ai/websockets/config"
	"github.com/cuttle-ai/websockets/log"
)

/*
 * this file contains the defintions of the rate limiter.
 * Basically the server cater the no. of requests at a given point of time as per specs.
 * When requests overflows it become very easy to scale if it is tracked.
 */

//RequestType is the type of the AppContext Request
type RequestType int

const (
	//Get is to get an app context
	Get RequestType = 0
	//Finished is to return an app context
	Finished RequestType = 1
	//CleanUp is to clean up the non-authenticated app context
	CleanUp RequestType = 2
	//Authenticate is to authenticate an unauthenticated request in system
	Authenticate RequestType = 3
)

//AppContextRequest is the request to get, return or try clean up app contexts
type AppContextRequest struct {
	//AppContext is the appcontext being requested
	AppContext *config.AppContext
	//Type is the type of request
	Type RequestType
	//Out is the ouput channel for get requests
	Out chan AppContextRequest
	//Exhausted flag states whether the app context exhausted
	Exhausted bool
	//Session is  the user session
	Session authConfig.Session
}

//AppContextRequestChan channel through which the app context routine takes requests from
var AppContextRequestChan = make(chan AppContextRequest)

//SendRequest is to send request to the channel. When this function used as go routines
//the blocking quenes can be solved
func SendRequest(ch chan AppContextRequest, req AppContextRequest) {
	ch <- req
}

//AppContext is the app context go routine running to
func AppContext(in chan AppContextRequest) {
	/*
	 * We will keep two maps for storing busy requests and free requests
	 * First we will generate the id pool and store it in
	 * We will start inifinite loop waiting for the requests
	 */
	//maps for storing the free and used requests
	freeMaps := make([]int, config.MaxRequests)
	unauthenticatedMap := make(map[int]time.Time, config.MaxRequests)
	authenticatedMap := make(map[int]time.Time, config.MaxRequests)

	//generate the request pool
	for i := 1; i <= config.MaxRequests; i++ {
		freeMaps = append(freeMaps, i)
	}

	//starting the infinite loop waiting for the requests
	for {
		req := <-in
		switch req.Type {
		case Get:
			//If it is a get request we will try to get get a app context from the store
			if len(freeMaps) == 0 {
				req.Exhausted = true
				go SendRequest(req.Out, req)
				return
			}
			id := freeMaps[0]
			freeMaps = freeMaps[1:]
			unauthenticatedMap[id] = time.Now()
			req.AppContext = config.NewAppContext(log.NewLogger(id))
			req.Exhausted = false
			go SendRequest(req.Out, req)
		case Finished:
			//we will return the request ids
			delete(authenticatedMap, req.AppContext.Log.GetID())
			freeMaps = append(freeMaps, req.AppContext.Log.GetID())
		case Authenticate:
			//we will return the request ids
			delete(unauthenticatedMap, req.AppContext.Log.GetID())
			authenticatedMap[req.AppContext.Log.GetID()] = time.Now()
		case CleanUp:
			//clean up the timed out requests
			n := time.Now()
			tot := config.IdleRequestTimeout
			maxLife := config.MaxRequestLife
			toBeAdded := []int{}
			for k, v := range unauthenticatedMap {
				if v.Add(tot).Before(n) {
					toBeAdded = append(toBeAdded, k)
					delete(unauthenticatedMap, k)
				}
			}
			for k, v := range authenticatedMap {
				if v.Add(maxLife).Before(n) {
					toBeAdded = append(toBeAdded, k)
					delete(authenticatedMap, k)
				}
			}
			freeMaps = append(freeMaps, toBeAdded...)
		}
	}
}

//CleanUpCheck is the cleanup check to be used as a go routine which periodically sends cleanup
//requests to the AppContext go routines
func CleanUpCheck(in chan AppContextRequest) {
	/*
	 * We will go into a infinte for loop
	 * Will send the requests of type clean up
	 */
	for {
		time.Sleep(config.RequestCleanUpCheck)
		go SendRequest(in, AppContextRequest{Type: CleanUp})
	}
}

func init() {
	go AppContext(AppContextRequestChan)
	go CleanUpCheck(AppContextRequestChan)
}
