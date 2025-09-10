package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
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

func (h *ChatHandler) MyConversation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user := userAny.(*userpb.UserPublicResponse)

	res, err := h.chatClient.GetConversationByUserId(ctx, &chatpb.GetByUserIdRequest{
		UserId: user.Id,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	common.JSON(c, http.StatusOK, "Lấy thông tin cuộc trò chuyện thành công", gin.H{
		"conversation": res,
	})
}

func (h *ChatHandler) ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.JSON(c, http.StatusInternalServerError, "Cập nhật kết nối từ HTTP -> WebSocket thất bại: "+err.Error(), nil)
		return
	}

	client := customWs.NewClient(h.hub, conn, c.Query("conversation_id"))

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
