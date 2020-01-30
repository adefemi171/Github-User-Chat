package main

import (
	"github.com/gorilla/websocket"
)

// client represents chatting with a single user.
type client struct {
	// socket is the websocket for the client
	socket *websocket.Conn
	// send is a channel on which messages are sent to the client
	send chan []byte
	// room is the place/room the client is chatting in.
	room *room
}

// read function for reading from the websocket
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// write function for writing to websocket
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
