// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package routes_test

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/cuttle-ai/websockets/routes"
	"github.com/cuttle-ai/websockets/routes/response"

	"github.com/cuttle-ai/websockets/config"
	"github.com/cuttle-ai/websockets/log"
)

/*
 * This file contains the examples required for the package documentation
 */

func ExampleInitRoutes() {
	//creating a new server mux
	m := http.NewServeMux()

	//created the default server
	s := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        m,
		ReadTimeout:    config.RequestRTimeout,
		WriteTimeout:   config.ResponseWTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//inited the routes
	routes.InitRoutes(m)

	//listen and serve to the server
	go func() {
		log.Info("Starting the User Subscription server at :" + config.Port)
		log.Error(s.ListenAndServe())
	}()

	//listening for syscalls
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt)
	sig := <-gracefulStop

	//gracefulling exiting when request comes in
	log.Info("Received the interrupt", sig)
	log.Info("Shutting down the server")
	err := s.Shutdown(context.Background())
	if err != nil {
		log.Error("Couldn't end the server gracefully")
	}
}

func ExampleAddRoutes() {
	//using add routes to create routes to handle requests
	routes.AddRoutes(routes.Route{
		Version: "v1",
		HandlerFunc: func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
			// rest of the implementation
		},
		Pattern: "/hi",
	})
}

func ExampleHandlerFunc() {
	//Example for creating a simple handler function
	f := func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		response.Write(res, response.Message{Message: "hi"})
	}
	routes.AddRoutes(routes.Route{
		Version:     "v1",
		HandlerFunc: f,
		Pattern:     "/hi",
	})
}

func ExampleHandlerFunc_context() {
	//Example for creating a simple handler function with context
	f := func(ctx context.Context, res http.ResponseWriter, req *http.Request) {

		//suppose there is a chan through a concurrent action happens
		tm := make(chan int)
		defer func() {
			//don't for get to close the channel
			close(tm)
		}()

		//kick start the concurrent action
		go func(ch chan int) {
			time.Sleep(1 * time.Second)
			ch <- 1
		}(tm)

		//wait for the results
		select {
		case <-tm:
			//we get the response
			response.Write(res, response.Message{Message: "hi"})
		case <-ctx.Done():
			//if timeout wins. Handle it gracefully.
			//No need to write the response
			log.Error("Timed out")
		}
	}
	routes.AddRoutes(routes.Route{
		Version:     "v1",
		HandlerFunc: f,
		Pattern:     "/hi",
	})
}
