package controller

import (
	"github.com/bulutcan99/go-websocket/internal/model"
	"github.com/bulutcan99/go-websocket/internal/platform/pubsub"
	"github.com/bulutcan99/go-websocket/internal/platform/repository"
	wsocket "github.com/bulutcan99/go-websocket/internal/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HubInterface interface {
	CreateRoom(roomName string) (*wsocket.Room, error)
}

type HubController struct {
	hub  *wsocket.Hub
	repo *repository.ChatRepo
	auth *AuthController
	pub  *pubsub.KafkaPublisher
	sub  *pubsub.KafkaSubscriber
}

func NewHubController(hub *wsocket.Hub, chatRepo *repository.ChatRepo, authCont *AuthController, pub *pubsub.KafkaPublisher, sub *pubsub.KafkaSubscriber) *HubController {
	return &HubController{
		hub:  hub,
		repo: chatRepo,
		auth: authCont,
		pub:  pub,
		sub:  sub,
	}
}

type CreateNewRoom struct {
	Name string `json:"room_name"`
}

func (h *HubController) CreateRoom(c *fiber.Ctx) error {
	var req CreateNewRoom
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Bad Request",
		})
	}

	roomId := uuid.New().String()
	h.hub.Rooms[roomId] = &wsocket.Room{
		UUID:    roomId,
		Name:    req.Name,
		Clients: make(map[string]*wsocket.Client),
	}

	return c.JSON(fiber.Map{
		"error": false,
		"Room":  req.Name,
	})
}

func (h *HubController) JoinRoom(c *websocket.Conn) {
	roomId := c.Params("roomID")
	userId := c.Query("userID")
	nickname := c.Query("nickname")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		c.WriteMessage(websocket.TextMessage, []byte("Room Not Found"))
		c.Close()
	}

	if _, ok := h.hub.Rooms[roomId].Clients[userId]; ok {
		c.WriteMessage(websocket.TextMessage, []byte("User Already Joined"))
		c.Close()
	}

	userUUID := uuid.MustParse(userId)
	roomUUID := uuid.MustParse(roomId)
	client := &wsocket.Client{
		Conn:     c,
		Message:  make(chan *wsocket.Message, 10),
		UUID:     userUUID,
		RoomUUID: roomUUID,
		Nickname: nickname,
	}

	h.hub.Rooms[roomId].Clients[userId] = client

	go client.WriteMessage()
	go client.ReadMessage(h.hub)
}

func (h *HubController) GetAllRooms(c *fiber.Ctx) error {
	var rooms []wsocket.Room
	for _, room := range h.hub.Rooms {
		rooms = append(rooms, *room)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"rooms": rooms,
	})
}

func (h *HubController) GetAvailableRooms(c *fiber.Ctx) error {
	var rooms []wsocket.Room
	for _, room := range h.hub.Rooms {
		if len(room.Clients) < 2 {
			rooms = append(rooms, *room)
		}
	}

	return c.JSON(fiber.Map{
		"error": false,
		"rooms": rooms,
	})
}

func (h *HubController) GetRoomClients(c *fiber.Ctx) error {
	roomId := c.Params("roomID")
	var clients []model.ClientResponse

	if _, ok := h.hub.Rooms[roomId]; !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "Room Not Found",
		})
	}

	for _, client := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, model.ClientResponse{
			UUID:     client.UUID,
			Nickname: client.Nickname,
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"clients": clients,
	})
}
