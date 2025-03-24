package pubsub

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type PubSub interface {
	message.Publisher
	message.Subscriber
}

func New(log *slog.Logger) PubSub {
	return gochannel.NewGoChannel(
		gochannel.Config{BlockPublishUntilSubscriberAck: false},
		watermill.NewSlogLogger(log),
	)
}
