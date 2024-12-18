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

func sanitizeFiles(input []string) []string {
	var output []string
	for _, f := range input {
		res := sanitize(f)
		if res != "" {
			output = append(output, res)
		}
	}

	return output
}

func (c *Client) Read(userID uint32) {
	defer c.Socket.Close()
	for {
		msg := &models.Message{}
		_, jsonMessage, err := c.Socket.ReadMessage()
		if err != nil {
			c.chatController.responder.LogError(fmt.Errorf("read message: %w", err), wc)
			return
		}

		err = easyjson.Unmarshal(jsonMessage, msg)
		if err != nil {
			c.chatController.responder.LogError(err, wc)
			return
		}
		msg.Content.Text = sanitize(msg.Content.Text)
		msg.Content.FilePath = sanitizeFiles(msg.Content.FilePath)
		msg.Content.StickerPath = sanitize(msg.Content.StickerPath)
		if msg.Content.Text == "" && msg.Content.StickerPath == "" && len(msg.Content.FilePath) == 0 {
			msg.Content.Text = "Я хотел отправить XSS"
		}
		msg.Sender = userID
		c.chatController.Messages <- msg
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	for msg := range c.Receive {
		msg.CreatedAt = time.Now()
		msg.Content.Text = sanitize(msg.Content.Text)
		msg.Content.FilePath = sanitizeFiles(msg.Content.FilePath)
		msg.Content.StickerPath = sanitize(msg.Content.StickerPath)
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
