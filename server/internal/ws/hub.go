package ws

type Room struct {
	ID      string             `json:"id"`
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
			if room, ok := h.Rooms[cl.RoomID]; ok {
				if _, exists := room.Clients[cl.ID]; !exists {
					room.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if room, ok := h.Rooms[cl.RoomID]; ok {
				if _, exists := room.Clients[cl.ID]; exists {
					if len(room.Clients) > 0 {
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}
					delete(room.Clients, cl.ID)
					close(cl.Message)
				}
			}
		case m := <-h.Broadcast:
			if room, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range room.Clients {
					cl.Message <- m
				}
			}
		}
	}
}