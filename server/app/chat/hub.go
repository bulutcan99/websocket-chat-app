package chat

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	UUID         uuid.UUID `json:"id"`
	Text         string    `json:"text"`
	CreatedAt    time.Time `json:"created_at"`
	RoomUUID     uuid.UUID `json:"room_id"`
	SenderUser   uuid.UUID `json:"sender_user"`
	ReceiverUser uuid.UUID `json:"receiver_user"`
}

type Client struct {
	UUID     uuid.UUID `json:"uuid"`
	Conn     *websocket.Conn
	Message  chan *Message
	RoomUUID uuid.UUID `json:"room_id"`
	Nickname string    `json:"nickname"`
}

type ChatRoom struct {
	UUID         uuid.UUID          `json:"id"`
	Name         string             `json:"name"`
	MessageCount int                `json:"message_count"`
	Clients      map[string]*Client `json:"clients"`
	CreatedAt    time.Time          `json:"created_at"`
	DeletedAt    *time.Time         `json:"deleted_at"`
}

type Hub struct {
	Rooms      map[string]*ChatRoom
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*ChatRoom),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}
