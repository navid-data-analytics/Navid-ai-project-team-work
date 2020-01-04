package aws

import (
	"context"

	"github.com/callstats-io/go-common/log"
)

// Publisher is a common interface for publishers
type Publisher interface {
	Publish(ctx context.Context, payload string) error
}

// LogPublisher implements the Publisher interface for offline/test usage
type LogPublisher struct {
	TopicARN string
}

// Assert that the LogPublisher conforms to the Publisher interface
var _ = Publisher(&LogPublisher{})

// NewLogPublisher builds and returns a new LogPublisher
func NewLogPublisher(topicArn string) *LogPublisher {
	return &LogPublisher{
		TopicARN: topicArn,
	}
}

// Publish logs a string formatted message within a context
func (p *LogPublisher) Publish(ctx context.Context, payload string) error {
	if payload == "" {
		logger(ctx).Error(LogEmptyPayload, log.Object(LogKeyOriginalPayload, payload))
		return ErrEmptyPayload
	}

	logger(ctx).Info(LogMessagePublished, log.String(LogKeyOriginalPayload, payload))
	return nil
}
