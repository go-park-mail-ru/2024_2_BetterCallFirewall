package chat

import (
	"github.com/gorilla/websocket"

	"github.com/2024_2_BetterCallFirewall/internal/chat/controller"
)

type Client struct {
	Socket         *websocket.Conn
	Receive        chan []byte
	chatController *controller.ChatController
}

func (c *Client) Read() {
	defer c.Socket.Close()
	for {
		_, jsonMessage, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		c.chatController.Messages <- jsonMessage
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Receive {
		err := c.Socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
