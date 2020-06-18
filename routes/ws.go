// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cuttle-ai/websockets/config"
	"github.com/cuttle-ai/websockets/models"
	"github.com/cuttle-ai/websockets/routes/response"
)

//WebSockets is the websockets connection handler
func WebSockets(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appCtx := ctx.Value(AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a websockets connection request")
	appCtx.WebSockets.ServeHTTP(res, req)
}

//SendNotification will send notification to connected websockets client of the user
func SendNotification(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	/*
	 * First we will get the app context
	 * Then we will parse the request payload
	 * Then we will get the web socket connection corresponding to the user
	 * Will write the response
	 * Then will send notification to the user
	 */
	//getting the app ctx
	appCtx := ctx.Value(AppContextKey).(*config.AppContext)
	appCtx.Log.Info("a request has come to send notification to the user", appCtx.Session.User.ID)

	//parse the request payload
	n := &models.Notification{}
	err := json.NewDecoder(req.Body).Decode(n)
	if err != nil {
		//bad request
		appCtx.Log.Error("error while parsing the notification", err.Error())
		response.WriteError(res, response.Error{Err: "Invalid Params " + err.Error()}, http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	//getting the user's websocket clients
	appCtxReq := AppContextRequest{
		Type:       FetchWs,
		Out:        make(chan AppContextRequest),
		AppContext: appCtx,
	}
	go SendRequest(AppContextRequestChan, appCtxReq)
	resCtx := <-appCtxReq.Out

	//sending response
	response.Write(res, response.Message{Message: "sending notitifications"})

	//sending notification to the user
	appCtx.Log.Info("sending notification event", n.Event, "to user", appCtx.Session.User.ID)
	for _, conn := range resCtx.WsConns {
		conn.Emit(n.Event, n.Payload)
	}
}

func init() {
	AddRoutes(Route{
		Version:     "v1",
		HandlerFunc: WebSockets,
		Pattern:     "/cuttle-websockets/",
	})
	AddRoutes(Route{
		Version:     "v1",
		HandlerFunc: SendNotification,
		Pattern:     "/notification/send",
	})
}
