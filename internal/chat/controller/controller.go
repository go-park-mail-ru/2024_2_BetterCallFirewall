package controller

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/2024_2_BetterCallFirewall/internal/chat"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestId string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestId string)

	ErrorInternal(w http.ResponseWriter, err error, requestId string)
	ErrorBadRequest(w http.ResponseWriter, err error, requestId string)
	LogError(err error, requestId string)
}

type ChatController struct {
	chatService chat.ChatService
	Messages    chan []byte
	responder   Responder
}

func NewChatController(service chat.ChatService, responder Responder) *ChatController {
	return &ChatController{
		chatService: service,
		Messages:    make(chan []byte),
		responder:   responder,
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var (
	upgrader    = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}
	mapUserConn = make(map[uint32]*Client)
)

func (cc *ChatController) SetConnection(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		cc.responder.LogError(myErr.ErrInvalidContext, "")
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		cc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	ctx := r.Context()

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		cc.responder.LogError(err, reqID)
		return
	}

	client := &Client{
		Socket:         socket,
		Receive:        make(chan []byte, messageBufferSize),
		chatController: cc,
	}
	mapUserConn[sess.UserID] = client
	defer func() {
		delete(mapUserConn, sess.UserID)
		close(client.Receive)
	}()
	go client.Write()
	go client.Read()
	cc.SendChatMsg(ctx, reqID, sess.UserID)
}

func (cc *ChatController) SendChatMsg(ctx context.Context, reqID string, userID uint32) {
	for jsonMsg := range cc.Messages {
		msg := &models.Message{}
		err := json.Unmarshal(jsonMsg, msg)
		if err != nil {
			cc.responder.LogError(err, reqID)
			return
		}
		msg.Sender = userID

		err = cc.chatService.SendNewMessage(ctx, msg.Receiver, msg.Sender, msg.Content)
		if err != nil {
			cc.responder.LogError(err, reqID)
			return
		}

		jsonForSend, err := json.Marshal(msg)
		if err != nil {
			cc.responder.LogError(err, reqID)
		}
		resConn, ok := mapUserConn[msg.Receiver]
		if ok {
			//resConn.Socket.ReadMessage()
			resConn.Receive <- jsonForSend
		}
	}
	return
}

func (cc *ChatController) GetAllChats(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok     = r.Context().Value("requestID").(string)
		lastTimeQuery = r.URL.Query().Get("lastIme")
		lastTime      time.Time
		err           error
	)

	if !ok {
		cc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	if lastTimeQuery == "" {
		lastTime = time.Now()
	} else {
		lastTime, err = time.Parse("2006-01-02 15:04:05.000", lastTimeQuery)
		if err != nil {
			cc.responder.ErrorBadRequest(w, myErr.ErrWrongDateFormat, reqID)
			return
		}
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		cc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	chats, err := cc.chatService.GetAllChats(r.Context(), sess.UserID, lastTime)
	if errors.Is(err, myErr.ErrNoMoreContent) {
		cc.responder.OutputNoMoreContentJSON(w, reqID)
		return
	}

	if err != nil {
		cc.responder.ErrorInternal(w, err, reqID)
		return
	}

	cc.responder.OutputJSON(w, chats, reqID)
}

func GetIdFromURL(r *http.Request) (uint32, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return 0, myErr.ErrEmptyId
	}

	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}
	if uid > math.MaxInt {
		return 0, myErr.ErrBigId
	}
	return uint32(uid), nil
}

func (cc *ChatController) GetChat(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok     = r.Context().Value("requestID").(string)
		lastTimeQuery = r.URL.Query().Get("lastIme")
		lastTime      time.Time
		err           error
	)

	if !ok {
		cc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	if lastTimeQuery == "" {
		lastTime = time.Now()
	} else {
		lastTime, err = time.Parse("2006-01-02 15:04:05.000", lastTimeQuery)
		if err != nil {
			cc.responder.ErrorBadRequest(w, myErr.ErrWrongDateFormat, reqID)
			return
		}
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		cc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	id, err := GetIdFromURL(r)
	if err != nil {
		cc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	messages, err := cc.chatService.GetChat(r.Context(), sess.UserID, id, lastTime)
	if errors.Is(err, myErr.ErrNoMoreContent) {
		cc.responder.OutputNoMoreContentJSON(w, reqID)
		return
	}

	if err != nil {
		cc.responder.ErrorInternal(w, err, reqID)
		return
	}

	cc.responder.OutputJSON(w, messages, reqID)
}
