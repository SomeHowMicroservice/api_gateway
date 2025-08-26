package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
)

type ChatContainer struct {
	Handler *handler.ChatHandler
}

func NewChatContainer(chatClient chatpb.ChatServiceClient) *ChatContainer {
	handler := handler.NewChatHandler(chatClient)
	return &ChatContainer{handler}
}
