package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SomeHowMicroservice/gateway/common"
	"github.com/SomeHowMicroservice/gateway/event"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/gin-contrib/sse"
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
	c.Header("Access-Control-Allow-Credentials", "true")
	origin := c.GetHeader("Origin")
	if origin != "" {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	userID := c.GetString("user_id")

	client := event.NewClient(userID)
	h.manager.Register <- client
	defer func() {
		h.manager.Unregister <- client
	}()

	sse.Encode(c.Writer, sse.Event{
		Event: "connected",
		Data:  gin.H{"message": "SSE connection established"},
	})
	c.Writer.Flush()

	clientGone := c.Request.Context().Done()

	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-client.Send:
			var msg common.SSEEvent
			if err := json.Unmarshal(message, &msg); err != nil {
				sse.Encode(c.Writer, sse.Event{
					Event: "error",
					Data:  gin.H{"message": fmt.Sprintf("parse failed: %v", err)},
				})
			} else {
				sse.Encode(c.Writer, sse.Event{
					Event: msg.Event,
					Data:  msg.Data,
				})
			}

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
