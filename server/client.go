package vcserver

import (
	"github.com/gorilla/websocket"
)

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			return
		}
	}
}