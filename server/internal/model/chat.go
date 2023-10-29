package model

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

const (
	MessageTypePresence = "presence"
	MessageTypeStatus   = "status"
	MessageTypeAuth     = "auth"
	MessageTypeMessage  = "message"

	MessageOffline = "0"
	MessageOnline  = "1"
)

// Message represents the message sent through websocket.
type Message struct {
	UUID         uuid.UUID           `json:"id"`
	Text         string              `json:"text"`
	Timestamp    timestamp.Timestamp `json:"timestamp"`
	SenderUser   uuid.UUID           `json:"sender_user"`
	ReceiverUser uuid.UUID           `json:"room_user"`
}

type Notification struct {
	UUID         uuid.UUID           `json:"id"`
	Notification string              `json:"notification"`
	Timestamp    timestamp.Timestamp `json:"timestamp"`
	SenderUser   uuid.UUID           `json:"sender_user"`
	ReceiverUser uuid.UUID           `json:"room_user"`
}
