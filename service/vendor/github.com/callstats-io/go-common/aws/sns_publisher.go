package aws

import (
	"context"

	"github.com/callstats-io/go-common/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SNSPublisher wraps Client and SNS topic ARN
type SNSPublisher struct {
	client   Client
	TopicARN string
}

// Assert that the SNSPublisher conforms to the Publisher interface
var _ = Publisher(&SNSPublisher{})

// NewSNSPublisher builds and returns a new SNSPublisher
func NewSNSPublisher(c Client, topicArn string) *SNSPublisher {
	return &SNSPublisher{
		client:   c,
		TopicARN: topicArn,
	}
}

// SNSClient returns a new SNS client based on the underlying AWS client session
func (p *SNSPublisher) SNSClient(ctx context.Context) (*sns.SNS, error) {
	sess, err := p.client.Session(ctx)
	if err != nil {
		return nil, err
	}
	return sns.New(sess), nil
}

// Publish sends a string formatted message to SNS topic
func (p *SNSPublisher) Publish(ctx context.Context, payload string) error {
	if payload == "" {
		logger(ctx).Error(LogInvalidSNSMessage, log.Object(LogKeyOriginalPayload, payload))
		return ErrInvalidSNSMessage
	}

	pi := &sns.PublishInput{
		Message:  aws.String(payload),
		TopicArn: aws.String(p.TopicARN),
	}

	snsClient, err := p.SNSClient(ctx)
	if err != nil {
		return err
	}

	output, err := snsClient.Publish(pi)
	if err != nil {
		logAWSErr(ctx, LogSNSPublishError, err)
		if output != nil {
			logger(ctx).Error(LogSNSPublishError, log.String(LogKeyAWSResponse, output.String()))
		}
		return err
	}

	logger(ctx).Info(LogSNSMessagePublished, log.String(LogKeyAWSResponse, output.String()))
	return nil
}
