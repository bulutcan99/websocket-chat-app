package wsocket

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"log"
	"time"
)

type Client struct {
	UUID     uuid.UUID `json:"uuid"`
	Conn     *websocket.Conn
	Message  chan *Message
	RoomUUID uuid.UUID `json:"room_id"`
	Nickname string    `json:"nickname"`
}

type Message struct {
	UUID           uuid.UUID `json:"id"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"created_at"`
	RoomUUID       uuid.UUID `json:"room_id"`
	SenderUser     string    `json:"sender_user"`
	ReceiverClient uuid.UUID `json:"receiver_client"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			UUID:           uuid.New(),
			Text:           string(m),
			CreatedAt:      time.Now(),
			RoomUUID:       c.RoomUUID,
			SenderUser:     c.Nickname,
			ReceiverClient: c.UUID,
		}

		hub.Broadcast <- msg
	}
}
