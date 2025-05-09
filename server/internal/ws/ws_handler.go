package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *fiber.Ctx) error {
	var req CreateRoomReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	return c.Status(fiber.StatusOK).JSON(req)
}

func (h *Handler) JoinRoom(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "websocket upgrade required"})
	}

	return websocket.New(func(conn *websocket.Conn) {
		roomID := c.Params("roomId")
		clientID := c.Query("userId")
		username := c.Query("username")

		client := &Client{
			Conn:     conn,
			Message:  make(chan *Message, 10),
			ID:       clientID,
			RoomID:   roomID,
			Username: username,
		}

		h.hub.Register <- client
		h.hub.Broadcast <- &Message{
			Content:  "A new user has joined the room",
			RoomID:   roomID,
			Username: username,
		}

		go client.writeMessage()
		client.readMessage(h.hub)
	})(c)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c *fiber.Ctx) error {
	rooms := make([]RoomRes, 0, len(h.hub.Rooms))
	for _, room := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}
	return c.JSON(rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *fiber.Ctx) error {
	roomID := c.Params("roomId")
	room, exists := h.hub.Rooms[roomID]
	if !exists {
		return c.JSON([]ClientRes{})
	}

	clients := make([]ClientRes, 0, len(room.Clients))
	for _, client := range room.Clients {
		clients = append(clients, ClientRes{
			ID:       client.ID,
			Username: client.Username,
		})
	}
	return c.JSON(clients)
}
