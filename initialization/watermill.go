package initialization

import (
	"fmt"

	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

type WatermillSubscriber struct {
	Subscriber message.Subscriber
}

func InitWatermill(cfg *config.AppConfig, logger watermill.LoggerAdapter, exchangeName string) (*WatermillSubscriber, error) {
	amqpConfig := amqp.NewDurablePubSubConfig(
		fmt.Sprintf("amqps://%s:%s@%s/%s",
			cfg.MessageQueue.RUser,
			cfg.MessageQueue.RPassword,
			cfg.MessageQueue.RHost,
			cfg.MessageQueue.RVhost,
		),
		nil,
	)

	amqpConfig.Exchange = amqp.ExchangeConfig{
		GenerateName: func(topic string) string {
			return exchangeName
		},
		Type:    "topic",
		Durable: true,
	}

	amqpConfig.Queue = amqp.QueueConfig{
		GenerateName: func(topic string) string {
			return topic
		},
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
	}

	amqpConfig.QueueBind = amqp.QueueBindConfig{
		GenerateRoutingKey: func(topic string) string {
			return topic
		},
	}

	amqpConfig.Consume.Qos = amqp.QosConfig{
		PrefetchCount: 5,
	}

	subscriber, err := amqp.NewSubscriber(amqpConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("tạo subscriber thất bại: %w", err)
	}

	return &WatermillSubscriber{
		subscriber,
	}, nil
}

func (w *WatermillSubscriber) Close() {
	_ = w.Subscriber.Close()
}
