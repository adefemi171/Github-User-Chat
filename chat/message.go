package main

import (
	"time"
)

// message that represents a single message
// a type to replace []byte slice
// This holds a time stamp of when the message was sent
type message struct {
	Name 		string
	Message 	string
	When		time.Time
	AvatarURL	string
}