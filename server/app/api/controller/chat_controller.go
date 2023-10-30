package controller

import (
	"github.com/bulutcan99/go-websocket/app/chat"
	"github.com/bulutcan99/go-websocket/internal/platform/pubsub"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type HubInterface interface {
	CreateRoom(roomName string) (*chat.ChatRoom, error)
}

type HubController struct {
	hub *chat.Hub
	pub pubsub.KafkaPublisher
	sub pubsub.KafkaSubscriber
}

func NewHubController(hub *chat.Hub, pub pubsub.KafkaPublisher, sub pubsub.KafkaSubscriber) HubController {
	return HubController{
		hub: hub,
		pub: pub,
		sub: sub,
	}
}

type CreateNewRoom struct {
	RoomName string `json:"room_name"`
}

func (h *HubController) CreateRoom(c *fiber.Ctx) error {
	var body CreateNewRoom
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	newRoomUUID := uuid.New()
	h.hub.Rooms[newRoomUUID.String()] = &chat.ChatRoom{
		UUID:         uuid.New(),
		Name:         body.RoomName,
		MessageCount: 0,
		Clients:      make(map[string]*chat.Client),
		CreatedAt:    time.Now(),
		DeletedAt:    nil,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "New room created",
	})
}

func (h *HubController) JoinRoom(c *websocket.Conn) error {
	roomUUID := c.Params("roomUUID")
	clientUUID := c.Query("userUUID")
	nickame := c.Query("nickname")

	if _, ok := h.hub.Rooms[roomUUID]; !ok {
		return c.WriteMessage(websocket.TextMessage, []byte("Room not found"))
		c.Close()
	} else {

	}

}
