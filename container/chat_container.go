package container

import (
	"github.com/SomeHowMicroservice/gateway/handler"
	chatpb "github.com/SomeHowMicroservice/gateway/protobuf/chat"
)

type ChatContainer struct {
	Handler *handler.ChatHandler
}

func NewChatContainer(chatClient chatpb.ChatServiceClient) *ChatContainer {
	handler := handler.NewChatHandler(chatClient)
	return &ChatContainer{handler}
}
