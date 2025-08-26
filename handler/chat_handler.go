package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatClient chatpb.ChatServiceClient
}

func NewChatHandler(chatClient chatpb.ChatServiceClient) *ChatHandler {
	return &ChatHandler{chatClient}
}

func (h *ChatHandler) TestConnect(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if _, err := h.chatClient.SendMessage(ctx, &chatpb.SendMessageRequest{
		Message: "Lồn",
	}); err != nil {
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Gửi tin nhắn thành công", nil)
}
