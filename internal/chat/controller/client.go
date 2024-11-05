package controller

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Client struct {
	Socket         *websocket.Conn
	Receive        chan []byte
	chatController *ChatController
}

func (c *Client) Read(userID uint32) {
	defer c.Socket.Close()
	for {
		msg := &models.Message{}
		_, jsonMessage, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}
		err = json.Unmarshal(jsonMessage, msg)
		if err != nil {
			return
		}
		msg.Sender = userID
		c.chatController.Messages <- msg
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
