package main

import (
	"errors"
)

// ErrNoAvatar is the error returned when the Avatar instance is unable to provide an avatar URL
var ErrNoAvatar = errors.New("chat: Unable to get an avatar URL.")

// Avatar represents types capable of representing user profile pictures
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client, or returns an error if something goes wrong
	GetAvatarURL(c *client) (string, error)

	// ErrNoAvatarURL is returned if the object is unable to get a URL for the specific specified client.
}
