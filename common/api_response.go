package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func JSON(c *gin.Context, statusCode int, message string, data interface{}) {
	if pb, ok := data.(proto.Message); ok {
		marshaler := protojson.MarshalOptions{
			EmitUnpopulated: true, 
			UseProtoNames: true,
		}
		jsonBytes, err := marshaler.Marshal(pb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ApiResponse{
				Message: "marshal protobuf thất bại",
				Data:    nil,
			})
			return
		}

		c.Data(statusCode, "application/json", []byte(`{"message":"`+message+`","data":`+string(jsonBytes)+`}`))
		return
	}
	
	c.JSON(statusCode, ApiResponse{
		Message: message,
		Data:    data,
	})
}
