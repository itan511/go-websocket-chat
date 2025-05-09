package ws

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				return
			}
			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("write error: %v", err)
				return
			}
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		if err := c.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read error: %v", err)
			}
			break
		}

		msg.RoomID = c.RoomID
		msg.Username = c.Username
		hub.Broadcast <- &msg
	}
}
