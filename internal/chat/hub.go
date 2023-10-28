package chat

import (
	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
	"time"
)

var (
	curr *websocket.Conn
	hubs *Hubs
)

type pardesData struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

type MessageData struct {
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
}

type Hubs struct {
	hubs map[string]*Hub
	run  chan *Hub
	stop chan *Hub
}

type Hub struct {
	chat       string
	current    chan *websocket.Conn
	clients    map[*websocket.Conn]string
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	running    chan bool
	isRunning  bool `default:"false"`
}

func getHubRun() *Hubs {
	return &Hubs{
		hubs: make(map[string]*Hub),
		run:  make(chan *Hub),
		stop: make(chan *Hub),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case connection := <-h.register:
			h.clients[connection] = connection.Params("conn_id")

		case message := <-h.broadcast:
			for connection := range h.clients {
				if curr != connection {
					if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
						h.unregister <- connection
						connection.WriteMessage(websocket.CloseMessage, []byte("Error Occured!"))
						connection.Close()
					}
				}
			}

		case connection := <-h.unregister:
			delete(h.clients, connection)
			if len(h.clients) == 0 {
				zap.S().Infof("initiating to stop hub")
				hubs.stop <- h
			}

		case curr := <-h.running:
			if !curr {
				return
			}

		case connection := <-h.current:
			curr = connection
		}
	}
}

func HubRunner() {
	hubs = getHubRun()
	for {
		select {
		case hub := <-hubs.run:
			zap.S().Infof("Starting hub %v", hub)
			go hub.Run()

		case hub := <-hubs.stop:
			zap.S().Infof("Stopping hub %v", hub)
			hub.running <- false
			delete(hubs.hubs, hub.chat)
		}
	}
}
