package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/event"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	manager    *event.Manager
	userClient userpb.UserServiceClient
}

func NewSSEHandler(manager *event.Manager, userClient userpb.UserServiceClient) *SSEHandler {
	return &SSEHandler{
		manager,
		userClient,
	}
}

func (h *SSEHandler) HandleSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	userID := c.Query("user_id")
	if userID == "" {
		common.JSON(c, http.StatusBadRequest, "user_id là bắt buộc", nil)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.userClient.CheckUserExistsById(ctx, &userpb.CheckUserExistsByIdRequest{
		Id: userID,
	})
	if common.HandleGrpcError(c, err) {
		return
	}

	if !res.Exists {
		common.JSON(c, http.StatusNotFound, common.ErrUserNotFound.Error(), nil)
		return
	}

	client := event.NewClient(userID)

	h.manager.Register <- client

	defer func() {
		h.manager.Unregister <- client
	}()

	c.Writer.Write([]byte("data: {\"event\":\"connected\",\"message\":\"SSE connection established\"}\n\n"))
	c.Writer.Flush()

	clientGone := c.Request.Context().Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-client.Send:
			fmt.Fprintf(c.Writer, "data: %s\n\n", string(message))
			c.Writer.Flush()

		case <-client.Done:
			return

		case <-clientGone:
			return

		case <-ticker.C:
			c.Writer.Write([]byte(": keep-alive\n\n"))
			c.Writer.Flush()
		}
	}
}
