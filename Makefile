proto-gen:
	@protoc --go_out=protobuf/auth --go-grpc_out=protobuf/auth \
	        proto/auth.proto
	@protoc --go_out=protobuf/user --go-grpc_out=protobuf/user \
	        proto/user.proto
	@protoc --go_out=protobuf/product --go-grpc_out=protobuf/product \
	        proto/product.proto
	@protoc --go_out=protobuf/post --go-grpc_out=protobuf/post \
	        proto/post.proto
	@protoc --go_out=protobuf/chat --go-grpc_out=protobuf/chat \
	        proto/chat.proto