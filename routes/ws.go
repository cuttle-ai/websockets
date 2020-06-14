// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes

import (
	"context"
	"net/http"

	"github.com/cuttle-ai/websockets/config"
)

//WebSockets is the websockets connection handler
func WebSockets(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	appCtx := ctx.Value(AppContextKey).(*config.AppContext)
	appCtx.Log.Info("Got a websockets connection request")
	appCtx.WebSockets.ServeHTTP(res, req)
}

func init() {
	AddRoutes(Route{
		Version:     "v1",
		HandlerFunc: WebSockets,
		Pattern:     "/cuttle-websockets/",
	})
}
