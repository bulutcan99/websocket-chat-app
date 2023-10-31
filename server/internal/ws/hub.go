package wsocket

type Room struct {
	UUID    string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Rooms[cl.RoomUUID.String()]; ok {
				r := h.Rooms[cl.RoomUUID.String()]

				if _, ok := r.Clients[cl.UUID.String()]; !ok {
					r.Clients[cl.UUID.String()] = cl
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomUUID.String()]; ok {
				if _, ok := h.Rooms[cl.RoomUUID.String()].Clients[cl.UUID.String()]; ok {
					if len(h.Rooms[cl.RoomUUID.String()].Clients) != 0 {
						h.Broadcast <- &Message{
							Text:       "user left the chat",
							RoomUUID:   cl.RoomUUID,
							SenderUser: cl.Nickname,
						}
					}

					delete(h.Rooms[cl.RoomUUID.String()].Clients, cl.UUID.String())
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomUUID.String()]; ok {

				for _, cl := range h.Rooms[m.RoomUUID.String()].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
