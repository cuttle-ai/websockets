// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package config will have necessary configuration for the application
package config

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cuttle-ai/websockets/version"

	"github.com/cuttle-ai/configs/config"
)

var (
	//Port in which the application is being served
	Port = "8078"
	//IntPort is the port converted into integer
	IntPort = 8078
	//RPCPort in which the application's rpc server is being served
	RPCPort = "8079"
	//RPCIntPort is the rpc port converted into integer
	RPCIntPort = 8079
	//RequestRTimeout of the api request body read timeout in milliseconds
	RequestRTimeout = time.Duration(2000 * time.Millisecond)
	//ResponseWTimeout of the api response write timeout in milliseconds
	ResponseWTimeout = time.Duration(10000 * time.Millisecond)
	//IdleRequestTimeout is the timeout after which unauthenticated requests must be disconnected of the system
	IdleRequestTimeout = time.Duration(10000 * time.Millisecond)
	//MaxRequestLife is the max request life time :- ie 4 hours is the default value
	MaxRequestLife = time.Duration(14400000)
	//MaxRequests is the maximum no. of requests catered at a given point of time
	MaxRequests = 1000
	//RequestCleanUpCheck is the time after which request cleanup check has to happen
	RequestCleanUpCheck = time.Duration(2 * time.Minute)
	//DiscoveryURL is the url of the discovery service
	DiscoveryURL = "127.0.0.1:8500"
	//DiscoveryToken is the token to communicate with discovery service
	DiscoveryToken = ""
	//ServiceDomain is the url on which the service will be available across the platform
	ServiceDomain = "127.0.0.1"
)

//SkipVault will skip the vault initialization if set true
var SkipVault bool

//IsTest indicates that the current runtime is for test
var IsTest bool

func init() {
	/*
	 * Based on the env variables will set the
	 *	* SkipVault
	 *  * IsTest
	 */
	sk := os.Getenv("SKIP_VAULT")
	if sk == "true" {
		SkipVault = true
	}
	iT := os.Getenv("IS_TEST")
	if iT == "true" {
		IsTest = true
	}
}

func init() {
	/*
	 * We will load the config from secrets management service
	 * Then we will set them as environment variables
	 */
	//getting the configuration
	log.Println("Getting the config values from vault")
	if SkipVault {
		return
	}
	v, err := config.NewVault()
	checkError(err)
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	configName := strings.ToLower(reg.ReplaceAllString(version.AppName, "-"))
	if IsTest {
		configName += "-test"
	}
	config, err := v.GetConfig(configName)
	checkError(err)

	//setting the configs as environment variables
	for k, v := range config {
		log.Println("Setting the secret from vault", k)
		os.Setenv(k, v)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	/*
	 * We will init the port
	 * We will init rpc port
	 * We will init the request timeout
	 * We will init the request body read timeout
	 * We will init the request body write timeout
	 * We will init the idle request timeout
	 * We will init the max request life
	 * We will init the max no. of requests
	 * We will init the request cleanup check
	 */
	//port
	if len(os.Getenv("PORT")) != 0 {
		//Assign the default port as 8078
		Port = os.Getenv("PORT")
		ip, err := strconv.Atoi(Port)
		if err != nil {
			//error whoile converting the port to integer
			log.Fatal("Error while converting the port to integer", err.Error())
		}
		IntPort = ip
	}

	//rpc port
	if len(os.Getenv("RPC_PORT")) != 0 {
		//Assign the default port as 8079
		RPCPort = os.Getenv("RPC_PORT")
		ip, err := strconv.Atoi(RPCPort)
		if err != nil {
			//error whoile converting the rpc port to integer
			log.Fatal("Error while converting the rpc port to integer", err.Error())
		}
		RPCIntPort = ip
	}

	//request body read timeout
	if len(os.Getenv("REQUEST_BODY_READ_TIMEOUT")) != 0 {
		//if successful convert timeout
		if t, err := strconv.ParseInt(os.Getenv("REQUEST_BODY_READ_TIMEOUT"), 10, 64); err == nil {
			RequestRTimeout = time.Duration(t * int64(time.Millisecond))
		}
	}

	//response write
	if len(os.Getenv("RESPOSE_WRITE_TIMEOUT")) != 0 {
		//if successful convert timeout
		if t, err := strconv.ParseInt(os.Getenv("RESPOSE_WRITE_TIMEOUT"), 10, 64); err == nil {
			ResponseWTimeout = time.Duration(t * int64(time.Millisecond))
		}
	}

	//idle request timeout
	if len(os.Getenv("IDLE_REQUEST_TIMEOUT")) != 0 {
		//if successful convert timeout
		if t, err := strconv.ParseInt(os.Getenv("IDLE_REQUEST_TIMEOUT"), 10, 64); err == nil {
			IdleRequestTimeout = time.Duration(t * int64(time.Millisecond))
		}
	}

	//max request life
	if len(os.Getenv("MAX_REQUEST_LIFE")) != 0 {
		//if successful convert life time
		if t, err := strconv.ParseInt(os.Getenv("MAX_REQUEST_LIFE"), 10, 64); err == nil {
			MaxRequestLife = time.Duration(t * int64(time.Millisecond))
		}
	}

	//max no. of requests
	if len(os.Getenv("MAX_REQUESTS")) != 0 {
		//if successful convert timeout
		if r, err := strconv.Atoi(os.Getenv("MAX_REQUESTS")); err == nil {
			MaxRequests = r
		}
	}

	//request cleanup check
	if len(os.Getenv("REQUEST_CLEAN_UP_CHECK")) != 0 {
		//if successful convert timeout
		if t, err := strconv.ParseInt(os.Getenv("REQUEST_CLEAN_UP_CHECK"), 10, 64); err == nil {
			RequestCleanUpCheck = time.Duration(t * int64(time.Minute))
		}
	}

	//discovery service url
	if len(os.Getenv("DISCOVERY_URL")) != 0 {
		DiscoveryURL = os.Getenv("DISCOVERY_URL")
	}

	//discovery service token
	if len(os.Getenv("DISCOVERY_TOKEN")) != 0 {
		DiscoveryToken = os.Getenv("DISCOVERY_TOKEN")
	}

	if len(DiscoveryToken) == 0 {
		log.Fatal("Token for discovery service is missing. Can't start the application without it")
	}

	//service domain
	if len(os.Getenv("SERVICE_DOMAIN")) != 0 {
		ServiceDomain = os.Getenv("SERVICE_DOMAIN")
	}
}

var (
	//PRODUCTION is the switch to turn on and off the Production environment.
	//1: On, 0: Off
	PRODUCTION = 0
)

func init() {
	/*
	 * Will init Production switch
	 */
	//Production
	if len(os.Getenv("PRODUCTION")) != 0 {
		//if successful convert production
		if t, err := strconv.Atoi(os.Getenv("PRODUCTION")); err == nil && (t == 1 || t == 0) {
			PRODUCTION = t
		}
	}
}
