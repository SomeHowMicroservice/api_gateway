package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatClient chatpb.ChatServiceClient
}

func NewChatHandler(chatClient chatpb.ChatServiceClient) *ChatHandler {
	return &ChatHandler{chatClient}
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
