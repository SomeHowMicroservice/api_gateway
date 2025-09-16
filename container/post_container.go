package container

import (
	"github.com/SomeHowMicroservice/gateway/handler"
	postpb "github.com/SomeHowMicroservice/gateway/protobuf/post"
)

type PostContainer struct {
	Handler *handler.PostHandler
}

func NewPostContainer(postClient postpb.PostServiceClient) *PostContainer {
	handler := handler.NewPostHandler(postClient)
	return &PostContainer{handler}
}
