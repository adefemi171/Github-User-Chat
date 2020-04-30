package main

import (
	"time"
	"github.com/gorilla/websocket"
)

// client represents chatting with a single user.
type client struct {
	// socket is the websocket for the client
	socket *websocket.Conn
	// send is a channel on which messages are sent to the client
	// send chan []byte
	send chan *message
	// room is the place/room the client is chatting in.
	room *room
	// userData holds information about the user
	userData map[string]interface{}
}

// read function for reading from the websocket
func (c *client) read() {
	defer c.socket.Close()
	for {
		// _, msg, err := c.socket.ReadMessage()
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarURL, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarURL.(string)
		}
		c.room.forward <- msg
	}
}

// write function for writing to websocket
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		// err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
			// return
		}
	}
}
