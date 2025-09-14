package handler

import (
	"net/http"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/SomeHowMicroservice/shm-be/gateway/socket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub *socket.Hub
}

func NewWSHandler(hub *socket.Hub) *WSHandler {
	return &WSHandler{hub}
}

func (h *WSHandler) HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.JSON(c, http.StatusInternalServerError, "Cập nhật kết nối từ HTTP -> WebSocket thất bại: "+err.Error(), nil)
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, common.ErrUnAuth.Error(), nil)
		return
	}

	user, _ := userAny.(*userpb.UserPublicResponse)
	convertedRole := convertMultiRolesToSingleRole(user.Roles)

	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		common.JSON(c, http.StatusBadRequest, "yêu cầu truyền conversation_id", nil)
		return
	}

	client := socket.NewClient(h.hub, conn, user.Id, convertedRole, conversationID)

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func convertMultiRolesToSingleRole(userRoles []string) string {
	if len(userRoles) > 1 {
		return "customer"
	}

	return "staff"
}
