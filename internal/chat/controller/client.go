package controller

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

const wc = "websocket"

type Client struct {
	Socket         *websocket.Conn
	Receive        chan *models.Message
	chatController *ChatController
}

func sanitize(input string) string {
	sanitizer := bluemonday.UGCPolicy()
	cleaned := sanitizer.Sanitize(input)
	return cleaned
}

func (c *Client) Read(userID uint32) {
	defer c.Socket.Close()
	for {
		msg := &models.Message{}
		_, jsonMessage, err := c.Socket.ReadMessage()
		clearJson := sanitize(string(jsonMessage))
		if err != nil {
			c.chatController.responder.LogError(fmt.Errorf("read message: %w", err), wc)
			return
		}
		err = easyjson.Unmarshal([]byte(clearJson), msg)

		if err != nil {
			c.chatController.responder.LogError(err, wc)
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
		jsonForSend, err := easyjson.Marshal(msg)
		if err != nil {
			c.chatController.responder.LogError(err, wc)
			return
		}
		err = c.Socket.WriteMessage(websocket.TextMessage, jsonForSend)
		if err != nil {
			c.chatController.responder.LogError(fmt.Errorf("write message: %w", err), wc)
			return
		}
	}
}
