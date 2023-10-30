package chat

import (
	"sync"
)

var MyServer = NewServer()

type Server struct {
	Clients    map[string]*Client
	mutex      *sync.Mutex
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		mutex:      &sync.Mutex{},
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}
