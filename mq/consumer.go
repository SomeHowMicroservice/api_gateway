package mq

import (
	"encoding/json"
	"fmt"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/event"
	"github.com/ThreeDotsLabs/watermill/message"
)

func RegisterProductImageUploadedConsumer(router *message.Router, subscriber message.Subscriber, sseManager *event.Manager) {
	router.AddConsumerHandler(
		"product_image_uploaded_handler",
		common.ProductUploadedTopic,
		subscriber,
		message.NoPublishHandlerFunc(func(msg *message.Message) error {
			return handleProductImageUploaded(msg, sseManager)
		}),
	)
}

func RegisterPostImageUploadedConsumer(router *message.Router, subscriber message.Subscriber, sseManager *event.Manager) {
	router.AddConsumerHandler(
		"post_image_uploaded_handler",
		common.PostUploadedTopic,
		subscriber,
		message.NoPublishHandlerFunc(func(msg *message.Message) error {
			return handlePostImageUploaded(msg, sseManager)
		}),
	)
}

func handleProductImageUploaded(msg *message.Message, sseManager *event.Manager) error {
	var eventMsg *common.ImageUploadedEvent
	if err := json.Unmarshal(msg.Payload, &eventMsg); err != nil {
		return fmt.Errorf("unmarshal json thất bại: %w", err)
	}

	event := &common.SSEEvent{
		Event: common.ProductImageUploaded,
		Data:  eventMsg,
	}

	if eventMsg.UserID == "" {
		return nil
	}

	sseManager.SendToUser(eventMsg.UserID, event)

	return nil
}

func handlePostImageUploaded(msg *message.Message, sseManager *event.Manager) error {
	var eventMsg *common.ImageUploadedEvent
	if err := json.Unmarshal(msg.Payload, &eventMsg); err != nil {
		return fmt.Errorf("unmarshal json thất bại: %w", err)
	}

	event := &common.SSEEvent{
		Event: common.PostImageUploaded,
		Data:  eventMsg,
	}

	if eventMsg.UserID == "" {
		return nil
	}

	sseManager.SendToUser(eventMsg.UserID, event)

	return nil
}
