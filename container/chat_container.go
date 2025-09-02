package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	"github.com/SomeHowMicroservice/shm-be/gateway/websocket"
)

type ChatContainer struct {
	Handler *handler.ChatHandler
}

func NewChatContainer(chatClient chatpb.ChatServiceClient, hub *websocket.Hub) *ChatContainer {
	handler := handler.NewChatHandler(chatClient, hub)
	return &ChatContainer{handler}
}
