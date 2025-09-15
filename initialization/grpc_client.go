package initialization

import (
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
	"google.golang.org/grpc/keepalive"
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

func InitClients(ca common.ClientAddresses) (*GRPCClients, error) {
	opts := dialOptions()
	authConn, err := grpc.NewClient(ca.AuthAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Auth Service: %w", err)
	}
	authClient := authpb.NewAuthServiceClient(authConn)

	userConn, err := grpc.NewClient(ca.UserAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Usẻ Service: %w", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	productConn, err := grpc.NewClient(ca.ProductAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Product Service: %w", err)
	}
	productClient := productpb.NewProductServiceClient(productConn)

	postConn, err := grpc.NewClient(ca.PostAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới Post Service: %w", err)
	}
	postClient := postpb.NewPostServiceClient(postConn)

	chatConn, err := grpc.NewClient(ca.ChatAddr, opts...)
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

func dialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{}],
				"retryPolicy": {
					"MaxAttempts": 4,
					"InitialBackoff": "0.1s",
					"MaxBackoff": "1s", 
					"BackoffMultiplier": 2.0,
					"RetryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]
				}
			}]
		}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
}
