package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	"github.com/SomeHowMicroservice/shm-be/gateway/socket"
)

type ChatContainer struct {
	Handler *handler.ChatHandler
}

func NewChatContainer(chatClient chatpb.ChatServiceClient, hub *socket.Hub) *ChatContainer {
	handler := handler.NewChatHandler(chatClient, hub)
	return &ChatContainer{handler}
}
