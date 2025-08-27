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
	authConn      *grpc.ClientConn
	userConn      *grpc.ClientConn
	productConn   *grpc.ClientConn
	postConn      *grpc.ClientConn
	chatConn      *grpc.ClientConn
}

func InitClients(ca *common.ClientAddresses) (*GRPCClients, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	authConn, err := grpc.DialContext(ctx, ca.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Auth Service: %w", err)
	}
	authClient := authpb.NewAuthServiceClient(authConn)

	userConn, err := grpc.DialContext(ctx, ca.UserAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Usẻ Service: %w", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	productConn, err := grpc.DialContext(ctx, ca.ProductAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Product Service: %w", err)
	}
	productClient := productpb.NewProductServiceClient(productConn)

	postConn, err := grpc.DialContext(ctx, ca.PostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Post Service: %w", err)
	}
	postClient := postpb.NewPostServiceClient(postConn)

	chatConn, err := grpc.DialContext(ctx, ca.ChatAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Chat Service: %w", err)
	}
	chatClient := chatpb.NewChatServiceClient(chatConn)

	return &GRPCClients{
		authClient,
		userClient,
		productClient,
		postClient,
		chatClient,
		authConn,
		userConn,
		productConn,
		postConn,
		chatConn,
	}, nil
}

func (g *GRPCClients) Close() {
	_ = g.authConn.Close()
	_ = g.userConn.Close()
	_ = g.productConn.Close()
	_ = g.postConn.Close()
	_ = g.chatConn.Close()
}
