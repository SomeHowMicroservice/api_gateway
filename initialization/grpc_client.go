package initialization

import (
	"context"
	"fmt"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	authpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/auth"
	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	postpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/post"
	productpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/product"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthClient    authpb.AuthServiceClient
	UserClient    userpb.UserServiceClient
	ProductClient productpb.ProductServiceClient
	PostClient    postpb.PostServiceClient
	ChatClient    chatpb.ChatServiceClient
}

func InitClients(ca *common.ClientAddresses) *GRPCClients {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	authConn, err := grpc.DialContext(ctx, ca.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		panic(fmt.Errorf("kết nối tới Auth Service thất bại: %w", err))
	}
	authClient := authpb.NewAuthServiceClient(authConn)

	userConn, err := grpc.DialContext(ctx, ca.UserAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		panic(fmt.Errorf("kết nối tới User Service thất bại: %w", err))
	}
	userClient := userpb.NewUserServiceClient(userConn)

	productConn, err := grpc.DialContext(ctx, ca.ProductAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		panic(fmt.Errorf("kết nối tới Product Service thất bại: %w", err))
	}
	productClient := productpb.NewProductServiceClient(productConn)

	postConn, err := grpc.DialContext(ctx, ca.PostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		panic(fmt.Errorf("kết nối tới Post Service thất bại: %w", err))
	}
	postClient := postpb.NewPostServiceClient(postConn)

	chatConn, err := grpc.DialContext(ctx, ca.ChatAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		panic(fmt.Errorf("kết nối tới Chat Service thất bại: %w", err))
	}
	chatClient := chatpb.NewChatServiceClient(chatConn)

	return &GRPCClients{
		authClient,
		userClient,
		productClient,
		postClient,
		chatClient,
	}
}
