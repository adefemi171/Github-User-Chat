package main

import (
	"log"
	"net/http"

	"github.com/adefemi171/Github-User-Chat/trace"
	"github.com/stretchr/objx"

	"github.com/gorilla/websocket"
)

// adding two channels
// first channel will add a client to the room
// second channel will remove a client from the room
type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients
	// forward chan []byte
	forward chan *message
	// join a channel for clients wishing to leave room
	join chan *client
	// leave is a channel for clients wishing to leave the room
	leave chan *client
	// clients holds all current clients in this room
	clients map[*client]bool
	// tracer will receive trace information of activity
	//in the room.
	tracer trace.Tracer
}

// newRoom chat function to create a new room
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

// using select statements to sychronize or modify shared memory
// run function that contains three select cases
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
			r.tracer.Trace("New Github User joined")
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("New Github User left")
		case msg := <-r.forward:
			r.tracer.Trace("Message Received: ", string(msg.Message))
			// forward message to all clients
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("-- Message Sent to Github User is ")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
