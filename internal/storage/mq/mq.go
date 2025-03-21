package mq

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type MessageQueue interface {
	message.Publisher
	message.Subscriber
}

func New(log *slog.Logger) MessageQueue {
	pubsub := gochannel.NewGoChannel(
		gochannel.Config{BlockPublishUntilSubscriberAck: false},
		watermill.NewSlogLogger(log),
	)

	return pubsub
}
