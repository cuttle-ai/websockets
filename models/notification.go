// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package models has the models required for the websockets service
package models

//Notification is the data translation object for sending notifications
type Notification struct {
	//Event is the event to be called
	Event string
	//Payload is the payload to be send with the notification
	Payload interface{}
}
