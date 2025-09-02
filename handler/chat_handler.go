package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	customWs "github.com/SomeHowMicroservice/shm-be/gateway/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	chatClient chatpb.ChatServiceClient
	hub        *customWs.Hub
}

func NewChatHandler(chatClient chatpb.ChatServiceClient, hub *customWs.Hub) *ChatHandler {
	return &ChatHandler{
		chatClient,
		hub,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *ChatHandler) TestConnect(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if _, err := h.chatClient.SendMessage(ctx, &chatpb.SendMessageRequest{
		Message: "Hello World!!!",
	}); err != nil {
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Gửi tin nhắn thành công", nil)
}

func (h *ChatHandler) ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.JSON(c, http.StatusInternalServerError, "Cập nhật kết nối từ HTTP -> WebSocket thất bại: "+err.Error(), nil)
		return
	}

	client := &customWs.Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: c.Query("user_id"),
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
