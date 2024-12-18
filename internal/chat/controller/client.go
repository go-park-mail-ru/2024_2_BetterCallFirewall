package controller

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Client struct {
	Socket         *websocket.Conn
	Receive        chan *models.Message
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
		msg.CreatedAt = time.Now()
		jsonForSend, err := json.Marshal(msg)
		if err != nil {
			return
		}
		err = c.Socket.WriteMessage(websocket.TextMessage, jsonForSend)
		if err != nil {
			return
		}
	}
}
